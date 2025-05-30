package domain

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user with such id not found")

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
