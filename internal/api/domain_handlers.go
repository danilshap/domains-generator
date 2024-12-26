package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/danilshap/domains-generator/internal/views/components/domains"
	"github.com/danilshap/domains-generator/internal/views/layouts"
	"github.com/go-chi/chi/v5"
)

const pageSize = 10

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/domains", http.StatusSeeOther)
}

func (s *Server) handleListDomains(w http.ResponseWriter, r *http.Request) {
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	offset := (page - 1) * pageSize

	domainsList, err := s.store.GetAllDomains(r.Context(), db.GetAllDomainsParams{
		Limit:  pageSize,
		Offset: int32(offset),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := s.store.GetDomainsCount(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalCount + pageSize - 1) / pageSize
	fmt.Println(totalPages)

	data := domains.ListData{
		Domains:     domainsList,
		CurrentPage: int32(page),
		TotalPages:  int32(totalPages),
		PageSize:    pageSize,
	}

	if r.Header.Get("HX-Request") == "true" {
		domains.TableWithPagination(data).Render(r.Context(), w)
		return
	}

	component := layouts.Base(domains.List(data))
	component.Render(r.Context(), w)
}

func (s *Server) handleNewDomainForm(w http.ResponseWriter, r *http.Request) {
	domains.Form().Render(r.Context(), w)
}

func isValidDomain(domain string) bool {
	pattern := `^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(domain)
}

func (s *Server) handleCreateDomain(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidDomain(r.FormValue("name")) {
		http.Error(w, "Invalid domain name format", http.StatusBadRequest)
		return
	}

	arg := db.CreateDomainParams{
		Name:     r.FormValue("name"),
		Provider: r.FormValue("provider"),
		Status:   1,
	}

	domain, err := s.store.CreateDomain(r.Context(), arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/domains/%d", domain.ID), http.StatusSeeOther)
}

func (s *Server) handleDeleteDomain(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.store.DeleteDomain(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	domainList, err := s.store.GetAllDomains(r.Context(), db.GetAllDomainsParams{
		Limit:  pageSize,
		Offset: 0,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := domains.ListData{
		Domains:     domainList,
		CurrentPage: 1,
		TotalPages:  1,
		PageSize:    pageSize,
	}
	domains.Table(data).Render(r.Context(), w)
}

func (s *Server) handleDomainDetails(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	offset := (page - 1) * mailboxPageSize

	domain, err := s.store.GetDomainByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mailboxes, err := s.store.GetMailboxesByDomainID(r.Context(), db.GetMailboxesByDomainIDParams{
		DomainID: domain.ID,
		Limit:    mailboxPageSize,
		Offset:   int32(offset),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := s.store.GetMailboxesCountByDomainID(r.Context(), domain.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalCount + mailboxPageSize - 1) / mailboxPageSize

	data := domains.DomainDetailsData{
		Domain:      domain,
		Mailboxes:   mailboxes,
		CurrentPage: int32(page),
		TotalPages:  int32(totalPages),
		PageSize:    mailboxPageSize,
	}

	if r.Header.Get("HX-Request") == "true" {
		domains.MailboxesTable(data).Render(r.Context(), w)
		return
	}

	component := layouts.Base(domains.DomainDetails(data))
	component.Render(r.Context(), w)
}

func (s *Server) handleEditDomainForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	domain, err := s.store.GetDomainByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	domains.EditForm(domain).Render(r.Context(), w)
}

func (s *Server) handleUpdateDomain(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	arg := db.UpdateDomainParams{
		ID:       int32(id),
		Name:     r.FormValue("name"),
		Provider: r.FormValue("provider"),
	}

	err = s.store.UpdateDomain(r.Context(), arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated list of domains
	domainsList, err := s.store.GetAllDomains(r.Context(), db.GetAllDomainsParams{
		Limit:  pageSize,
		Offset: 0,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := s.store.GetDomainsCount(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	data := domains.ListData{
		Domains:     domainsList,
		CurrentPage: 1,
		TotalPages:  int32(totalPages),
		PageSize:    pageSize,
	}

	domains.List(data).Render(r.Context(), w)
}

func (s *Server) handleStatusForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	domain, err := s.store.GetDomainByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	domains.StatusForm(domain).Render(r.Context(), w)
}

func (s *Server) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := strconv.Atoi(r.FormValue("status"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.store.UpdateDomainAndMailboxesStatus(r.Context(), int32(id), int32(status))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated list of domains
	domainsList, err := s.store.GetAllDomains(r.Context(), db.GetAllDomainsParams{
		Limit:  pageSize,
		Offset: 0,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := s.store.GetDomainsCount(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	data := domains.ListData{
		Domains:     domainsList,
		CurrentPage: 1,
		TotalPages:  int32(totalPages),
		PageSize:    pageSize,
	}

	domains.List(data).Render(r.Context(), w)
}
