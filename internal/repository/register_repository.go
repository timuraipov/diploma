package repository

import (
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/storage/db"
)

type registerRepository struct {
	database db.DB
}

func NewRegisterRepository(db db.DB) domain.DomainMock {
	return &registerRepository{
		database: db,
	}
}
