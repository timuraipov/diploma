package domain

import "context"

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Balance struct {
	ID        int64   `json:"id" omitempty:"true"`
	UserID    string  `json:"user_id"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
type BalanceRepository interface {
	GetBalance(ctx context.Context, userID int64) (Balance, error)
	UpdateBalance(ctx context.Context, userID int64, amount float64) error
	Withdraw(ctx context.Context, withdraw Withdraw) error
	Withdrawals(ctx context.Context, userID int64) ([]Withdraw, error)
}
type BalanceUseCase interface {
	GetBalance(ctx context.Context, userID int64) (Balance, error)
	UpdateBalance(ctx context.Context, userID int64, amount float64) error
	Withdraw(ctx context.Context, withdraw Withdraw) error
	Withdrawals(ctx context.Context, userID int64) ([]Withdraw, error)
}
