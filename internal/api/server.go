package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danilshap/domains-generator/internal/auth"
	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/danilshap/domains-generator/internal/middleware"
	"github.com/danilshap/domains-generator/pkg/config"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	store      db.Store
	router     *chi.Mux
	tokenMaker *auth.JWTMaker
}

func NewServer(store db.Store, cfg *config.Config) (*Server, error) {
	tokenMaker, err := auth.NewJWTMaker(cfg.TokenSynnetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		router:     chi.NewRouter(),
		tokenMaker: tokenMaker,
	}

	// Middleware
	server.router.Use(chimiddleware.Logger)
	server.router.Use(chimiddleware.Recoverer)
	server.router.Use(middleware.AuthMiddleware(server.tokenMaker))

	// Routes
	server.router.Get("/login", server.handleLoginPage)
	server.router.Get("/register", server.handleRegisterPage)
	server.router.Post("/login", server.handleLogin)
	server.router.Post("/register", server.handleRegister)
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
		r.Put("/{id}/status", s.handleUpdateMailboxStatus)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func getCurrentUser(r *http.Request) (*auth.Payload, error) {
	payload, ok := r.Context().Value(middleware.UserContextKey).(*auth.Payload)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	return payload, nil
}
