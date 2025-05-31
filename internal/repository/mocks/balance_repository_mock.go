package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/timuraipov/diploma/internal/domain"
)

type MockBalanceRepository struct {
	mock.Mock
	balances    map[int64]domain.Balance
	withdrawals map[string]domain.Withdraw
}

func NewMockBalanceRepository() MockBalanceRepository {
	balances := make(map[int64]domain.Balance)
	withdrawals := make(map[string]domain.Withdraw)
	return MockBalanceRepository{balances: balances, withdrawals: withdrawals}
}

func (b *MockBalanceRepository) GetBalance(ctx context.Context, userID int64) (domain.Balance, error) {
	for _, balance := range b.balances {
		if balance.UserID == userID {
			return balance, nil
		}
	}
	return domain.Balance{}, domain.ErrUserNotFound
}

func (b *MockBalanceRepository) UpdateBalance(ctx context.Context, userID int64, accrual float64) error {
	var balanceObj domain.Balance
	for _, balance := range b.balances {
		if balance.UserID == userID {
			balanceObj = balance
			break
		}
	}
	if balanceObj.ID == 0 {
		balanceObj.Balance = accrual
		balanceObj.CreatedAt = time.Now()
		balanceObj.ID = 777
		balanceObj.UserID = userID
		b.balances[balanceObj.UserID] = balanceObj
	} else {
		balanceObj.Balance += accrual
	}
	return nil
}
func (b *MockBalanceRepository) Withdraw(ctx context.Context, withdraw domain.Withdraw) error {
	_, ok := b.withdrawals[withdraw.ID]
	if ok {
		return domain.ErrOrderAlreadyProcessedByAnotherUser
	}
	balance, err := b.GetBalance(ctx, withdraw.UserID)
	if err != nil || balance.Balance < withdraw.Sum {
		return domain.ErrInsufficientFunds
	}
	b.withdrawals[withdraw.ID] = withdraw
	balance.Balance -= withdraw.Sum

	return nil
}
func (b *MockBalanceRepository) Withdrawals(ctx context.Context, userID int64) ([]domain.Withdraw, error) {
	withdrawals := make([]domain.Withdraw, 0)
	for _, withdraw := range b.withdrawals {
		if withdraw.UserID == userID {
			withdrawals = append(withdrawals, withdraw)
		}
	}
	return withdrawals, nil
}
