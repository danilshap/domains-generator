package api

import (
	"fmt"
	"net/http"
	"strconv"

	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/danilshap/domains-generator/internal/views/components/mailboxes"
	"github.com/danilshap/domains-generator/internal/views/layouts"
	"github.com/go-chi/chi/v5"
)

const mailboxPageSize = 10

func (s *Server) handleListMailboxes(w http.ResponseWriter, r *http.Request) {
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	var statusFilters []int32
	if status := r.URL.Query().Get("status"); status != "" {
		if s, err := strconv.Atoi(status); err == nil {
			statusFilters = append(statusFilters, int32(s))
		}
	}

	var domainFilters []int32
	if domain := r.URL.Query().Get("domain"); domain != "" {
		if d, err := strconv.Atoi(domain); err == nil {
			domainFilters = append(domainFilters, int32(d))
		}
	}

	search := r.URL.Query().Get("search")
	offset := (page - 1) * mailboxPageSize

	mailboxList, err := s.store.GetMailboxesWithFilters(r.Context(), db.GetMailboxesWithFiltersParams{
		Column1: statusFilters,
		Column2: domainFilters,
		Column3: search,
		Limit:   mailboxPageSize,
		Offset:  int32(offset),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats, err := s.store.GetMailboxesStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := s.store.GetMailboxesCount(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalCount + mailboxPageSize - 1) / mailboxPageSize

	data := mailboxes.ListData{
		Mailboxes:   mailboxList,
		CurrentPage: int32(page),
		TotalPages:  int32(totalPages),
		PageSize:    mailboxPageSize,
		Stats:       stats,
		Filters: struct {
			Status  []int32
			Domains []int32
			Search  string
		}{
			Status:  statusFilters,
			Domains: domainFilters,
			Search:  search,
		},
	}

	if r.Header.Get("HX-Request") == "true" {
		mailboxes.TableWithPagination(data).Render(r.Context(), w)
		return
	}

	component := layouts.Base(mailboxes.List(data))
	component.Render(r.Context(), w)
}

func (s *Server) handleNewMailboxForm(w http.ResponseWriter, r *http.Request) {
	domains, err := s.store.GetAllDomains(r.Context(), db.GetAllDomainsParams{
		Limit:  100, // Assuming there aren't too many domains
		Offset: 0,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mailboxes.Form(domains).Render(r.Context(), w)
}

func (s *Server) handleCreateMailbox(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	domainID, err := strconv.Atoi(r.FormValue("domain_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check that domain is active
	domain, err := s.store.GetDomainByID(r.Context(), int32(domainID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if domain.Status != 1 {
		http.Error(w, "Cannot create mailbox for inactive domain", http.StatusBadRequest)
		return
	}

	arg := db.CreateMailboxParams{
		Address:  r.FormValue("address"),
		Password: r.FormValue("password"),
		DomainID: int32(domainID),
		Status:   1, // Active by default
	}

	mailbox, err := s.store.CreateMailbox(r.Context(), arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/mailboxes/%d", mailbox.ID), http.StatusSeeOther)
}

func (s *Server) handleMailboxDetails(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id <= 0 {
		http.Error(w, "Invalid mailbox ID", http.StatusBadRequest)
		return
	}

	mailbox, err := s.store.GetMailboxByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	domain, err := s.store.GetDomainByID(r.Context(), mailbox.DomainID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := mailboxes.DetailsData{
		Mailbox: mailbox,
		Domain:  domain,
	}

	component := layouts.Base(mailboxes.Details(data))
	component.Render(r.Context(), w)
}

func (s *Server) handleEditMailboxForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mailbox, err := s.store.GetMailboxByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	domains, err := s.store.GetAllDomains(r.Context(), db.GetAllDomainsParams{
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mailboxes.EditForm(mailbox, domains).Render(r.Context(), w)
}

func (s *Server) handleUpdateMailbox(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	domainID, err := strconv.Atoi(r.FormValue("domain_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	arg := db.UpdateMailboxParams{
		ID:       int32(id),
		Address:  r.FormValue("address"),
		DomainID: int32(domainID),
	}

	err = s.store.UpdateMailbox(r.Context(), arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/mailboxes/%d", id), http.StatusSeeOther)
}

func (s *Server) handleDeleteMailbox(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.store.DeleteMailbox(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
