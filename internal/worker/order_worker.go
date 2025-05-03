package worker

import (
	"context"
	"net/http"
	"time"

	"github.com/timuraipov/diploma/internal/client"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/pkg/logging"
	"go.uber.org/zap"
)

type OrderWorker struct {
	l                 *logging.ZapLogger
	orderRepository   domain.OrderRepository
	balanceRepository domain.BalanceRepository
	client            client.AccrualClient
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewOrderWorker(l *logging.ZapLogger, orderRepository domain.OrderRepository, balanceRepository domain.BalanceRepository, client client.AccrualClient) *OrderWorker {
	ctx, cancel := context.WithCancel(context.Background()) // todo pass parent context
	return &OrderWorker{
		l:                 l,
		orderRepository:   orderRepository,
		balanceRepository: balanceRepository,
		client:            client,
		ctx:               ctx,
		cancel:            cancel,
	}
}
func (o *OrderWorker) Start() {
	const numJobs = 5
	jobs := make(chan domain.Order, numJobs)

	go o.StartTicker(o.ctx, o.orderRepository, jobs)
	for w := 1; w <= 3; w++ {
		go o.Worker(o.ctx, jobs)
	}
	o.l.InfoCtx(o.ctx, "Starting order worker")

}

func (o *OrderWorker) StartTicker(ctx context.Context, repository domain.OrderRepository, jobs chan domain.Order) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			orders, err := repository.GetUnhandledOrders(ctx)
			if err != nil {
				o.l.ErrorCtx(ctx, "Error getting orders- "+err.Error())
				continue
			}
			for _, order := range orders {
				o.l.InfoCtx(ctx, "Processing order", zap.Any("order", order))
				jobs <- order
			}
		case <-ctx.Done():
			return
		}
	}
}
func (o *OrderWorker) Worker(ctx context.Context, jobs chan domain.Order) {
	for { // TODO mutex on client
		select {
		case <-ctx.Done():
			return
		case order := <-jobs:
			o.l.InfoCtx(ctx, "Processing order", zap.Any("order", order))
			result, statusCode, err := o.client.GetOrder(order.ID)
			if err != nil {
				o.l.ErrorCtx(ctx, "Error getting order from client- "+err.Error())
				continue
			}
			if statusCode == http.StatusOK {
				order.Accrual = result.Amount
				order.Status = result.Status
				err := o.orderRepository.UpdateOrder(ctx, order)
				if err != nil {
					o.l.ErrorCtx(ctx, "Error updating order- "+err.Error())
					continue
				}
				o.l.InfoCtx(ctx, "Order updated successfully", zap.Any("order", order))
				if order.Status == "PROCESSED" {
					err := o.balanceRepository.UpdateBalance(ctx, order.UserId, order.Accrual)
					if err != nil {
						o.l.ErrorCtx(ctx, "Error updating balance- "+err.Error())
						continue
					}
					o.l.InfoCtx(ctx, "Balance updated successfully", zap.Any("userId", order.UserId), zap.Any("accrual", order.Accrual))
				}
			} else if statusCode == http.StatusNotFound {
				o.l.ErrorCtx(ctx, "Order not found", zap.Any("order", order))
			} else if statusCode == http.StatusInternalServerError {
				o.l.ErrorCtx(ctx, "Internal server error", zap.Any("order", order))
			}
			if err != nil {
				o.l.ErrorCtx(ctx, "Error saving order- "+err.Error())
			}
			continue
		}

	}
}
