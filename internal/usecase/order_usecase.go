package usecase

import (
	"context"
	"time"

	"github.com/timuraipov/diploma/internal/client"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/pkg/logging"
)

type OrderUseCase struct {
	l               *logging.ZapLogger
	orderRepository domain.OrderRepository
	accrualClient   client.AccrualClient
	contextTimeout  time.Duration
}

func NewOrderUseCase(l *logging.ZapLogger, ur domain.OrderRepository, client client.AccrualClient, timeout time.Duration) domain.OrderUseCase {
	return &OrderUseCase{
		l:               l,
		orderRepository: ur,
		contextTimeout:  timeout,
	}
}

func (o *OrderUseCase) Save(ctx context.Context, order domain.Order) error {
	err := o.orderRepository.Save(ctx, order)

	return err
}
func (o *OrderUseCase) GetAll(ctx context.Context, userId int64) ([]domain.Order, error) {
	orders, err := o.orderRepository.GetAll(ctx, userId)
	return orders, err
}
