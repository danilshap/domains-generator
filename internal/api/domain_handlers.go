package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/danilshap/domains-generator/internal/models/view"
	"github.com/danilshap/domains-generator/internal/services"
	"github.com/danilshap/domains-generator/internal/views/components/domains"
	"github.com/danilshap/domains-generator/internal/views/layouts"
	"github.com/go-chi/chi/v5"
)

const pageSize = 10
const mailboxPageSize = 10

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/domains", http.StatusSeeOther)
}

func (s *Server) handleListDomains(w http.ResponseWriter, r *http.Request) {
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

	offset := (page - 1) * pageSize

	domainsList, err := s.store.GetAllDomains(r.Context(), db.GetAllDomainsParams{
		Limit:  pageSize,
		Offset: int32(offset),
		UserID: user.UserID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := s.store.GetDomainsCount(r.Context(), user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	data := domains.ListData{
		Domains:     domainsList,
		CurrentPage: int32(page),
		TotalPages:  int32(totalPages),
		PageSize:    pageSize,
	}

	if r.Header.Get("HX-Request") == "true" {
		domains.List(data).Render(r.Context(), w)
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
	user, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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
		UserID:   user.UserID,
	}

	domain, err := s.store.CreateDomain(r.Context(), arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.domainService.CreateDomain(r.Context(), domain.Name); err != nil {
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

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Trigger", `{"showMessage": "Domain updated successfully"}`)
		w.Header().Set("HX-Redirect", "/domains")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/domains", http.StatusSeeOther)
}

func (s *Server) handleDomainDetails(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	domain, err := s.store.GetDomainByID(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Domain not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	emailAccounts, err := s.store.GetMailboxesByDomainID(r.Context(), db.GetMailboxesByDomainIDParams{
		DomainID: domain.ID,
		Limit:    int32(mailboxPageSize),
		Offset:   int32((page - 1) * mailboxPageSize),
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

	totalPages := (totalCount + int64(mailboxPageSize) - 1) / int64(mailboxPageSize)

	mailboxViews := make([]view.MailboxView, len(emailAccounts))
	for i, m := range emailAccounts {
		mailboxViews[i] = view.ToMailboxViewFromDomain(m)
	}

	data := domains.DetailsData{
		Domain:      domain,
		Mailboxes:   mailboxViews,
		CurrentPage: int32(page),
		TotalPages:  int32(totalPages),
		PageSize:    int32(mailboxPageSize),
	}

	if r.Header.Get("HX-Request") == "true" {
		domains.Details(data).Render(r.Context(), w)
		return
	}

	component := layouts.Base(domains.Details(data))
	component.Render(r.Context(), w)
}

func (s *Server) handleUpdateDomainStatus(w http.ResponseWriter, r *http.Request) {
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

	status, err := strconv.Atoi(r.FormValue("status"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if status != 1 && status != 2 {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	err = s.store.UpdateDomainAndMailboxesStatus(r.Context(), int32(id), int32(status))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status == 2 {
		go func() {
			domain, err := s.store.GetDomainByID(context.Background(), int32(id))
			if err != nil {
				log.Printf("Failed to get domain for notification: %v\n", err)
				return
			}

			notification, err := s.store.CreateNotification(context.Background(), db.CreateNotificationParams{
				UserID:  user.UserID,
				Title:   "Domain disabled",
				Message: fmt.Sprintf("Domain %s status updated to %d", domain.Name, status),
				Type:    db.NotificationTypeInfo,
			})
			if err != nil {
				log.Printf("Failed to create notification: %v\n", err)
				return
			}

			log.Printf("Sending WebSocket notification for user %d\n", user.UserID)
			s.wsService.SendNotification(user.UserID, services.WSNotification{
				Title:   notification.Title,
				Message: notification.Message,
				Type:    string(notification.Type),
			})
		}()
	}

	if r.Header.Get("HX-Request") == "true" {
		s.handleDomainDetails(w, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/domains/%d", id), http.StatusSeeOther)
}

func (s *Server) handleBulkMailboxesForm(w http.ResponseWriter, r *http.Request) {
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

	domains.BulkMailboxesForm(domain).Render(r.Context(), w)
}
