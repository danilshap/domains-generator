package middleware

import (
	"context"
	"net/http"

	"github.com/danilshap/domains-generator/internal/auth"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(tokenMaker *auth.JWTMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for login and register pages
			if r.URL.Path == "/login" || r.URL.Path == "/register" {
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := r.Cookie("token")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			payload, err := tokenMaker.VerifyToken(cookie.Value)
			if err != nil {
				http.SetCookie(w, &http.Cookie{
					Name:   "token",
					Value:  "",
					Path:   "/",
					MaxAge: -1,
				})
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
