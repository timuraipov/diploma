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

// UpdateBalance implements domain.BalanceUseCase.
func (b *balanceUseCase) UpdateBalance(ctx context.Context, userID int64, amount float64) error {
	err := b.balanceRepository.UpdateBalance(ctx, userID, amount)
	return err
}

func (b *balanceUseCase) GetBalance(ctx context.Context, userID int64) (domain.BalanceResponse, error) {
	balance, err := b.balanceRepository.GetBalance(ctx, userID)
	if err != nil {
		return domain.BalanceResponse{}, err
	}
	withdrawals, err := b.balanceRepository.Withdrawals(ctx, userID)
	if err != nil {
		return domain.BalanceResponse{}, err
	}
	var withdrawn float64
	for _, withraw := range withdrawals {
		withdrawn += withraw.Sum
	}
	return domain.BalanceResponse{
		Current:   balance.Balance,
		Withdrawn: withdrawn,
	}, nil
}
func (b *balanceUseCase) Withdraw(ctx context.Context, withdraw domain.Withdraw) error {
	err := b.balanceRepository.Withdraw(ctx, withdraw) //return 422 if order does not exist
	return err
}
func (b *balanceUseCase) Withdrawals(ctx context.Context, userID int64) ([]domain.WithdrawResponse, error) {
	withdrawals, err := b.balanceRepository.Withdrawals(ctx, userID)
	var withdrawResponse []domain.WithdrawResponse
	for _, withdraw := range withdrawals {
		withdrawRFC := &domain.WithdrawResponse{
			Sum:         withdraw.Sum,
			ID:          withdraw.ID,
			ProcessedAt: withdraw.ProcessedAt.Format(time.RFC3339),
		}
		withdrawResponse = append(withdrawResponse, *withdrawRFC)
	}
	return withdrawResponse, err
}
