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

type RegisterResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
type RegisterUsecase interface {
	Create(ctx context.Context, user RegisterRequest) (User, error)
	GetByLogin(ctx context.Context, login string) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	CreateAccessToken(user *User, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user *User, secret string, expiry int) (refreshToken string, err error)
}
