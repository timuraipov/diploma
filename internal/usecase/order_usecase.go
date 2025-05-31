package usecase

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"unicode"

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
	ok := o.luhnCheck(order.ID)
	o.l.InfoCtx(ctx, "Checking order ID", zap.Any("orderID", order.ID), zap.Bool("isOk", ok))
	if !ok {
		return domain.ErrIncorrectOrderIDFormat
	}
	err := o.orderRepository.Save(ctx, order)

	return err
}
func (o *OrderUseCase) GetAll(ctx context.Context, userID int64) ([]domain.Order, error) {
	orders, err := o.orderRepository.GetAll(ctx, userID)
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
			o.l.InfoCtx(ctx, "Order is processed, updating balance", zap.Any("userID", order.UserID), zap.Any("accrual", order.Accrual))
			err := o.balanceUseCase.UpdateBalance(ctx, order.UserID, order.Accrual)
			if err != nil {
				o.l.ErrorCtx(ctx, "Error updating balance- "+err.Error())
				return statusCode, err
			}
			o.l.InfoCtx(ctx, "Balance updated successfully", zap.Any("userID", order.UserID), zap.Any("accrual", order.Accrual))
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
func (o *OrderUseCase) luhnCheck(number string) bool {
	var sum int
	alt := false

	// Обрабатываем цифры справа налево
	for i := len(number) - 1; i >= 0; i-- {
		r := rune(number[i])
		if !unicode.IsDigit(r) {
			return false // если в строке есть нецифры
		}
		n, _ := strconv.Atoi(string(r))

		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}

	return sum%10 == 0
}
