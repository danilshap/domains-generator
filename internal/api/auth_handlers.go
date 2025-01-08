package api

import (
	"net/http"
	"time"

	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/danilshap/domains-generator/internal/views/components/auth"
	"github.com/danilshap/domains-generator/internal/views/layouts"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := loginRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	user, err := s.store.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := s.tokenMaker.CreateToken(user.ID, user.Email, 24*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(24 * time.Hour.Seconds()),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := registerRequest{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if len(req.Username) < 4 {
		http.Error(w, "Username must be at least 4 characters long", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := s.store.CreateUser(r.Context(), db.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: string(hashedPassword),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := s.tokenMaker.CreateToken(user.ID, user.Email, 24*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(24 * time.Hour.Seconds()),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	component := layouts.Auth(auth.Login())
	component.Render(r.Context(), w)
}

func (s *Server) handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	component := layouts.Auth(auth.Register())
	component.Render(r.Context(), w)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
