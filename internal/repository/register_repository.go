package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/storage/db"
)

type registerRepository struct {
	database db.DB
}

func NewRegisterRepository(db db.DB) *registerRepository {
	return &registerRepository{
		database: db,
	}
}
func (r *registerRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO "user" (login, password) VALUES (@login, @password)`
	args := pgx.NamedArgs{
		"login":    user.Login,
		"password": user.Password,
	}
	_, err := r.database.Pool.Exec(ctx, query, args)
	return err
}

func (r *registerRepository) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	return domain.User{}, nil
}
func (r *registerRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	return domain.User{}, nil
}
