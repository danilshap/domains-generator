package api

import (
	"fmt"
	"net/http"
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
		domains.Table(data.Domains).Render(r.Context(), w)
		return
	}

	component := layouts.Base(domains.List(data))
	component.Render(r.Context(), w)
}

func (s *Server) handleNewDomainForm(w http.ResponseWriter, r *http.Request) {
	domains.Form().Render(r.Context(), w)
}

func (s *Server) handleCreateDomain(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	domains.Table(data.Domains).Render(r.Context(), w)
}

func (s *Server) handleDomainDetails(w http.ResponseWriter, r *http.Request) {
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

	component := layouts.Base(domains.DomainDetails(domain))
	component.Render(r.Context(), w)
}
