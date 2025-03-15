package repository

import (
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/storage/db"
)

type loginRepository struct {
	database db.DB
}

func NewLoginRepository(db db.DB) domain.DomainMock {
	return &loginRepository{
		database: db,
	}
}
