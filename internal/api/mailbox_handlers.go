package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/danilshap/domains-generator/internal/models/view"
	"github.com/danilshap/domains-generator/internal/views/components/mailboxes"
	"github.com/danilshap/domains-generator/internal/views/layouts"
	"github.com/danilshap/domains-generator/pkg/utils"
	"github.com/go-chi/chi/v5"
)

func (s *Server) handleListMailboxes(w http.ResponseWriter, r *http.Request) {
	user, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	dimainId := int32(0)
	domainIDFilter := []int32{}
	if domainIDStr := r.URL.Query().Get("domain_id"); domainIDStr != "" {
		domainIDInt, err := strconv.Atoi(domainIDStr)
		if err == nil && domainIDInt > 0 {
			dimainId = int32(domainIDInt)
			domainIDFilter = append(domainIDFilter, dimainId)
		}
	}

	offset := (page - 1) * mailboxPageSize

	mailboxList, err := s.store.GetMailboxesWithFilters(r.Context(), db.GetMailboxesWithFiltersParams{
		StatusFilter: []int32{},      // status filters
		DomainFilter: domainIDFilter, // domain filters
		SearchQuery:  "",             // search
		UserID:       user.UserID,
		PageLimit:    mailboxPageSize,
		PageOffset:   int32(offset),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := s.store.GetMailboxesCount(r.Context(), db.GetMailboxesCountParams{
		UserID:       user.UserID,
		DomainFilter: domainIDFilter,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalCount + mailboxPageSize - 1) / mailboxPageSize

	mailboxViews := make([]view.MailboxView, len(mailboxList))
	for i, m := range mailboxList {
		mailboxViews[i] = view.ToMailboxViewFromFilter(m)
	}

	data := view.MailboxListData{
		Mailboxes:   mailboxViews,
		CurrentPage: int32(page),
		TotalPages:  int32(totalPages),
		PageSize:    mailboxPageSize,
		DomainID:    dimainId,
	}

	if r.Header.Get("HX-Request") == "true" {
		mailboxes.List(data).Render(r.Context(), w)
		return
	}

	component := layouts.Base(mailboxes.Page(data))
	component.Render(r.Context(), w)
}

func (s *Server) handleNewMailboxForm(w http.ResponseWriter, r *http.Request) {
	user, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	domainID := int32(0)
	if domainIDStr := r.URL.Query().Get("domain_id"); domainIDStr != "" {
		if domainIDInt, err := strconv.Atoi(domainIDStr); err == nil && domainIDInt > 0 {
			domainID = int32(domainIDInt)
		}
	}

	domains := []db.GetAllDomainsRow{}
	if domainID > 0 {
		domains, err = s.store.GetAllDomains(r.Context(), db.GetAllDomainsParams{
			Limit:  100,
			Offset: 0,
			UserID: user.UserID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	mailboxes.Form(domains, domainID).Render(r.Context(), w)
}

func (s *Server) handleCreateMailbox(w http.ResponseWriter, r *http.Request) {
	user, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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

	password := r.FormValue("password")
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	arg := db.CreateMailboxParams{
		Address:  r.FormValue("address"),
		Password: hashedPassword,
		DomainID: int32(domainID),
		UserID:   user.UserID,
		Status:   1,
	}

	// First try to create mailbox through provider
	err = s.mailboxService.CreateMailboxWithPassword(r.Context(), arg.Address, domain.Name, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Then save to database
	mailbox, err := s.store.CreateMailbox(r.Context(), arg)
	if err != nil {
		// Try to cleanup provider mailbox if database insert fails
		if delErr := s.mailboxService.DeleteMailbox(r.Context(), arg.Address); delErr != nil {
			log.Printf("Failed to cleanup mailbox after db error: %v", delErr)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Trigger", `{"showMessage": "Mailbox created successfully"}`)
		w.Header().Set("HX-Redirect", "/mailboxes")
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/mailboxes/%d", mailbox.ID), http.StatusSeeOther)
	}
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

	mailboxes.EditForm(mailbox).Render(r.Context(), w)
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

	newPassword := r.FormValue("password")
	if newPassword == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get existing mailbox
	mailbox, err := s.store.GetMailboxByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get domain for provider update
	domain, err := s.store.GetDomainByID(r.Context(), mailbox.DomainID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update password in provider
	if err := s.mailboxService.UpdatePassword(r.Context(), mailbox.Address, domain.Name, newPassword); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update password in database
	err = s.store.UpdateMailboxPassword(r.Context(), db.UpdateMailboxPasswordParams{
		ID:       int32(id),
		Password: hashedPassword,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Trigger", `{"showMessage": "Mailbox status updated successfully"}`)
		w.Header().Set("HX-Redirect", fmt.Sprintf("/mailboxes/%d", id))
		w.WriteHeader(http.StatusOK)
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

	// Get mailbox to know its address
	mailbox, err := s.store.GetMailboxByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// First delete from provider
	if err := s.mailboxService.DeleteMailbox(r.Context(), mailbox.Address); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Then delete from database
	err = s.store.DeleteMailbox(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleUpdateMailboxStatus(w http.ResponseWriter, r *http.Request) {
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

	// Validate status
	if status != 1 && status != 2 {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	// Get mailbox to know its address and domain
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

	// If deactivating, delete from provider
	if status == 2 {
		if err := s.mailboxService.DeleteMailbox(r.Context(), mailbox.Address); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if status == 1 { // If activating, recreate in provider
		if err := s.mailboxService.CreateMailboxWithPassword(r.Context(), mailbox.Address, domain.Name, mailbox.Password); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = s.store.UpdateMailboxesStatusByID(r.Context(), db.UpdateMailboxesStatusByIDParams{
		ID:     int32(id),
		Status: int32(status),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Trigger", `{"showMessage": "Mailbox status updated successfully"}`)
		w.Header().Set("HX-Redirect", fmt.Sprintf("/mailboxes/%d", id))
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/mailboxes/%d", id), http.StatusSeeOther)
}

func (s *Server) handleCreateBulkMailboxes(w http.ResponseWriter, r *http.Request) {
	user, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	prefix := r.FormValue("prefix")
	password := r.FormValue("password")

	if password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create 100 mailboxes
	addresses, err := s.mailboxService.CreateBulkMailboxes(r.Context(), prefix, domain.Name, password, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save mailboxes to database
	for _, address := range addresses {
		_, err := s.store.CreateMailbox(r.Context(), db.CreateMailboxParams{
			Address:  address,
			Password: hashedPassword,
			DomainID: domain.ID,
			UserID:   user.UserID,
			Status:   1,
		})
		if err != nil {
			// Log error but continue
			log.Printf("Failed to create mailbox %s: %v", address, err)
		}
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Trigger", `{"showMessage": "Bulk mailboxes created successfully"}`)
		w.Header().Set("HX-Redirect", fmt.Sprintf("/domains/%d", domain.ID))
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/domains/%d", domain.ID), http.StatusSeeOther)
}
