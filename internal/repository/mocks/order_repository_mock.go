package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/timuraipov/diploma/internal/domain"
)

type MockOrderRepository struct {
	mock.Mock
	db map[string]domain.Order
}

func NewMockOrderRepository() MockOrderRepository {
	db := make(map[string]domain.Order)
	return MockOrderRepository{db: db}
}

func (o *MockOrderRepository) Save(ctx context.Context, order domain.Order) error {
	orderFound, ok := o.db[order.ID]
	if ok {
		if orderFound.UserID == order.UserID {
			return domain.ErrOrderAlreadyInProcessing
		}
		return domain.ErrOrderAlreadyProcessedByAnotherUser
	}
	o.db[order.ID] = order
	return nil
}

func (o *MockOrderRepository) GetAll(ctx context.Context, userID int64) ([]domain.Order, error) {
	orders := make([]domain.Order, 0)
	for _, order := range o.db {
		if order.UserID == userID {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (o *MockOrderRepository) GetUnhandledOrders(ctx context.Context) ([]domain.Order, error) {
	orders := make([]domain.Order, 0)
	for _, order := range o.db {
		if order.Status == domain.REGISTERED || order.Status == domain.PROCESSING {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (o *MockOrderRepository) UpdateOrder(ctx context.Context, order domain.Order) error {
	orderFound, ok := o.db[order.ID]
	if ok {
		if orderFound.UserID == order.UserID {
			o.db[order.ID] = order
			return nil
		} else {
			return domain.ErrOrderAlreadyProcessedByAnotherUser
		}
	}
	return domain.ErrOrderNotFound
}
