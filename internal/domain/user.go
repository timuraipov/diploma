package domain

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound             = errors.New("user with such id not found")
	ErrIncorrectLoginOrPassword = errors.New("login request has incorrect login or password")
	ErrUserIDNotFound           = errors.New("cannot get userId from context")
)

type AuthSecret struct {
	AccessTokenExpiryHour  int
	RefreshTokenExpiryHour int
	AccessTokenSecret      string
	RefreshTokenSecret     string
}
type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByLogin(ctx context.Context, login string) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
}
