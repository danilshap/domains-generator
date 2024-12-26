package api

import (
	"net/http"

	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	store  db.Store
	router *chi.Mux
}

func NewServer(store db.Store) (*Server, error) {
	server := &Server{
		store:  store,
		router: chi.NewRouter(),
	}

	// Middleware
	server.router.Use(middleware.Logger)
	server.router.Use(middleware.Recoverer)

	// Routes
	server.setupRoutes()

	return server, nil
}

func (s *Server) setupRoutes() {
	s.router.Get("/", s.handleHome)
	s.router.Route("/domains", func(r chi.Router) {
		r.Get("/", s.handleListDomains)
		r.Get("/new", s.handleNewDomainForm)
		r.Post("/", s.handleCreateDomain)
		r.Get("/{id}", s.handleDomainDetails)
		r.Get("/{id}/edit", s.handleEditDomainForm)
		r.Get("/{id}/status", s.handleStatusForm)
		r.Put("/{id}", s.handleUpdateDomain)
		r.Put("/{id}/status", s.handleUpdateStatus)
		r.Delete("/{id}", s.handleDeleteDomain)
	})
	s.router.Route("/mailboxes", func(r chi.Router) {
		r.Get("/", s.handleListMailboxes)
		r.Get("/new", s.handleNewMailboxForm)
		r.Post("/", s.handleCreateMailbox)
		r.Get("/{id}", s.handleMailboxDetails)
		r.Get("/{id}/edit", s.handleEditMailboxForm)
		r.Put("/{id}", s.handleUpdateMailbox)
		r.Delete("/{id}", s.handleDeleteMailbox)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
