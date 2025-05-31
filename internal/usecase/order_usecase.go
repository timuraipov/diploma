package usecase

import (
	"context"
	"net/http"
	"strconv"
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
		balanceUseCase:  balanceUseCase,
		contextTimeout:  timeout,
	}
}

func (o *OrderUseCase) Save(ctx context.Context, order domain.Order) error {
	_, err := o.luhnCheckDigit(order.ID)
	if err != nil {
		return domain.ErrIncorrectOrderIDFormat
	}
	err = o.orderRepository.Save(ctx, order)

	return err
}
func (o *OrderUseCase) GetAll(ctx context.Context, userId int64) ([]domain.Order, error) {
	orders, err := o.orderRepository.GetAll(ctx, userId)
	return orders, err
}
func (o *OrderUseCase) Accrual(ctx context.Context, order domain.Order) (int, error) {

	o.l.InfoCtx(ctx, "Processing order", zap.Any("order-------", order))
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
			o.l.InfoCtx(ctx, "Order is processed, updating balance", zap.Any("userId", order.UserId), zap.Any("accrual", order.Accrual))
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
	o.l.InfoCtx(ctx, "Getting unhandled orders", zap.Any("orders", orders))
	if err != nil {
		o.l.ErrorCtx(ctx, "Error getting unhandled orders- "+err.Error())
		return nil, err
	}
	return orders, nil
}
func (o *OrderUseCase) luhnCheckDigit(s string) (int, error) {
	number, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	checkNumber := o.luhnChecksum(number)

	if checkNumber == 0 {
		return 0, nil
	}
	return 10 - checkNumber, nil
}

func (o *OrderUseCase) luhnChecksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
