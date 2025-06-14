package domain

import (
	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const UserIDHeader contextKey = "x-user-id"

type JwtCustomClaims struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	ID int64 `json:"id"`
	jwt.StandardClaims
}
