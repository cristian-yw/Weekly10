package models

import "github.com/golang-jwt/jwt/v5"

type JWTClaim struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
