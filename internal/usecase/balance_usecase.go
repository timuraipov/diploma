package usecase

import (
	"context"
	"time"

	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/pkg/logging"
)

type balanceUseCase struct {
	l                 *logging.ZapLogger
	balanceRepository domain.BalanceRepository
	contextTimeout    time.Duration
}

func NewBalanceUseCase(l *logging.ZapLogger, br domain.BalanceRepository, timeout time.Duration) domain.BalanceUseCase {
	return &balanceUseCase{
		l:                 l,
		balanceRepository: br,
		contextTimeout:    timeout,
	}
}

// GetBalance implements domain.BalanceUseCase.
func (b *balanceUseCase) GetBalance(ctx context.Context, userID int64) (domain.Balance, error) {
	balance, err := b.balanceRepository.GetBalance(ctx, userID)
	return balance, err
}

// UpdateBalance implements domain.BalanceUseCase.
func (b *balanceUseCase) UpdateBalance(ctx context.Context, userID int64, amount float64) error {
	err := b.balanceRepository.UpdateBalance(ctx, userID, amount)
	return err
}

// Withdraw implements domain.BalanceUseCase.
func (b *balanceUseCase) Withdraw(ctx context.Context, withdraw domain.Withdraw) error {
	err := b.balanceRepository.Withdraw(ctx, withdraw)

	return err
}

// Withdrawals implements domain.BalanceUseCase.
func (b *balanceUseCase) Withdrawals(ctx context.Context, userID int64) ([]domain.Withdraw, error) {
	withdrawals, err := b.balanceRepository.Withdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}
