package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/danilshap/domains-generator/internal/auth"
	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/danilshap/domains-generator/internal/middleware"
	"github.com/danilshap/domains-generator/internal/services"
	"github.com/danilshap/domains-generator/internal/services/domain"
	"github.com/danilshap/domains-generator/internal/services/mailbox"
	"github.com/danilshap/domains-generator/pkg/config"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

type Server struct {
	store          db.Store
	router         *chi.Mux
	tokenMaker     *auth.JWTMaker
	domainService  *domain.Service
	mailboxService *mailbox.Service
	wsService      *services.WebSocketService
}

func NewServer(store db.Store, cfg *config.Config) (*Server, error) {
	tokenMaker, err := auth.NewJWTMaker(cfg.TokenSynnetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	mockProvider := domain.NewMockProvider()
	domainService := domain.NewService(mockProvider)

	mockMailboxProvider := mailbox.NewMockProvider()
	mailboxService := mailbox.NewService(mockMailboxProvider)

	wsService := services.NewWebSocketService()

	server := &Server{
		store:          store,
		router:         chi.NewRouter(),
		tokenMaker:     tokenMaker,
		domainService:  domainService,
		mailboxService: mailboxService,
		wsService:      wsService,
	}

	// Middleware
	server.router.Use(chimiddleware.Logger)
	server.router.Use(chimiddleware.Recoverer)
	server.router.Use(middleware.AuthMiddleware(server.tokenMaker))

	// Static files
	server.router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Routes
	server.router.Get("/login", server.handleLoginPage)
	server.router.Get("/register", server.handleRegisterPage)
	server.router.Post("/login", server.handleLogin)
	server.router.Post("/register", server.handleRegister)
	server.router.Post("/logout", server.handleLogout)
	server.setupRoutes()

	return server, nil
}

func (s *Server) setupRoutes() {
	s.router.Get("/", s.handleHome)
	s.router.Get("/ws", s.handleWebSocket)
	s.router.Route("/domains", func(r chi.Router) {
		r.Get("/", s.handleListDomains)
		r.Get("/new", s.handleNewDomainForm)
		r.Post("/", s.handleCreateDomain)
		r.Get("/{id}", s.handleDomainDetails)
		r.Delete("/{id}", s.handleDeleteDomain)
		r.Put("/{id}/status", s.handleUpdateDomainStatus)
		r.Get("/{id}/bulk-mailboxes", s.handleBulkMailboxesForm)
		r.Post("/{id}/bulk-mailboxes", s.handleCreateBulkMailboxes)
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

	s.router.Route("/notifications", func(r chi.Router) {
		r.Get("/", s.handleListNotifications)
		r.Post("/{id}/mark-read", s.handleMarkNotificationRead)
		r.Post("/mark-all-read", s.handleMarkAllNotificationsRead)
		r.Get("/unread-count", s.handleUnreadNotificationsCount)
	})
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	user, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}

	s.wsService.AddConnection(user.UserID, conn)
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
