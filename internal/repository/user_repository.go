package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/storage/db"
)

type userRepository struct {
	database db.DB
}

func NewUserRepository(db db.DB) *userRepository {
	return &userRepository{
		database: db,
	}
}
func (u *userRepository) Create(ctx context.Context, user *domain.User) error {
	const stmt = `INSERT INTO "user" (login, password) VALUES (@login, @password) returning (id)`
	args := pgx.NamedArgs{
		"login":    user.Login,
		"password": user.Password,
	}
	err := u.database.Pool.QueryRow(ctx, stmt, args).Scan(&user.ID)
	return err
}

func (u *userRepository) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	user := domain.User{}
	const stmt = `SELECT id, login, password from "user" WHERE login=@login`
	args := pgx.NamedArgs{
		"login": login,
	}
	err := u.database.Pool.QueryRow(ctx, stmt, args).Scan(&user.ID, &user.Login, &user.Password)

	return user, err
}
func (u *userRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	user := domain.User{}
	const stmt = `SELECT id, login, password from "user" WHERE id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}
	err := u.database.Pool.QueryRow(ctx, stmt, args).Scan(&user.ID, &user.Login, &user.Password)

	return user, err
}
