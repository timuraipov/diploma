package domain

import (
	"context"
	"errors"
)

var (
	ErrUserAlreadyRegistered = errors.New("user already registered")
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
type RegisterUsecase interface {
	Create(ctx context.Context, user RegisterRequest) (AuthResponse, error)
	GetByLogin(ctx context.Context, login string) (User, error)
}
