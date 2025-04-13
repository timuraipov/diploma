package repository

import (
	"context"

	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/storage/db"
)

type balanceRepository struct {
	database db.DB
}

func NewBalanceRepository(db db.DB) *balanceRepository {
	return &balanceRepository{
		database: db,
	}
}

func (b *balanceRepository) GetBalance(ctx context.Context, balance domain.Balance) error {
	return nil
}
func (b *balanceRepository) UpdateBalance(ctx context.Context, balance domain.Balance) error {
	return nil
}

func (b *balanceRepository) Withdraw(ctx context.Context, withdraw domain.Withdraw) error {
	return nil
}
func (b *balanceRepository) Withdrawals(ctx context.Context, userID int64) ([]domain.Withdraw, error) {
	return nil, nil
}
