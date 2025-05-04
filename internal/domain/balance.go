package domain

import (
	"context"
	"errors"
	"time"
)

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

var (
	InsufficientFundsError = errors.New("insufficient funds for this operation")
)

type Balance struct {
	ID        int64     `json:"id" omitempty:"true"`
	UserID    int64     `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type BalanceRepository interface {
	GetBalance(ctx context.Context, userID int64) (Balance, error)
	UpdateBalance(ctx context.Context, userID int64, accrual float64) error
	Withdraw(ctx context.Context, withdraw Withdraw) error
	Withdrawals(ctx context.Context, userID int64) ([]Withdraw, error)
}
type BalanceUseCase interface {
	GetBalance(ctx context.Context, userID int64) (BalanceResponse, error)
	UpdateBalance(ctx context.Context, userID int64, amount float64) error
	Withdraw(ctx context.Context, withdraw Withdraw) error
	Withdrawals(ctx context.Context, userID int64) ([]WithdrawResponse, error)
}
