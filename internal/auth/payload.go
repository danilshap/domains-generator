package auth

import "github.com/golang-jwt/jwt/v5"

type Payload struct {
	UserID int32  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
