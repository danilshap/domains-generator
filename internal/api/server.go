package api

import (
	"net/http"

	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	store  *db.Store
	router *chi.Mux
}

func NewServer(store *db.Store) (*Server, error) {
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
		r.Delete("/{id}", s.handleDeleteDomain)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
