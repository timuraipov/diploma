package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/timuraipov/diploma/internal/client"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/pkg/logging"
	"go.uber.org/zap"
)

type OrderUseCase struct {
	l               *logging.ZapLogger
	orderRepository domain.OrderRepository
	client          client.AccrualClient
	balanceUseCase  domain.BalanceUseCase
	contextTimeout  time.Duration
}

func NewOrderUseCase(l *logging.ZapLogger, ur domain.OrderRepository, client client.AccrualClient, balanceUseCase domain.BalanceUseCase, timeout time.Duration) domain.OrderUseCase {
	return &OrderUseCase{
		l:               l,
		orderRepository: ur,
		client:          client,
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
func (o *OrderUseCase) Accrual(ctx context.Context, order domain.Order) (int, error) {

	o.l.InfoCtx(ctx, "Processing order", zap.Any("order", order))
	result, statusCode, err := o.client.GetOrder(order.ID)
	if err != nil {
		o.l.ErrorCtx(ctx, "Error getting order from client- "+err.Error())
		return statusCode, err
	}
	o.l.InfoCtx(ctx, "Received order from client", zap.Any("result", result), zap.Any("statusCode", statusCode))
	if statusCode == http.StatusOK {
		order.Accrual = result.Amount
		order.Status = result.Status
		err := o.orderRepository.UpdateOrder(ctx, order)
		if err != nil {
			o.l.ErrorCtx(ctx, "Error updating order- "+err.Error())
			return statusCode, err
		}
		o.l.InfoCtx(ctx, "Order updated successfully", zap.Any("order", order))
		if order.Status == "PROCESSED" {
			err := o.balanceUseCase.UpdateBalance(ctx, order.UserId, order.Accrual)
			if err != nil {
				o.l.ErrorCtx(ctx, "Error updating balance- "+err.Error())
				return statusCode, err
			}
			o.l.InfoCtx(ctx, "Balance updated successfully", zap.Any("userId", order.UserId), zap.Any("accrual", order.Accrual))
		}
	}
	return statusCode, err
}

func (o *OrderUseCase) GetUnhandledOrders(ctx context.Context) ([]domain.Order, error) {
	orders, err := o.orderRepository.GetUnhandledOrders(ctx)
	if err != nil {
		o.l.ErrorCtx(ctx, "Error getting unhandled orders- "+err.Error())
		return nil, err
	}
	return orders, nil
}
