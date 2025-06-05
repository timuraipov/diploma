package worker

import (
	"context"
	"net/http"
	"time"

	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/client"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/repository"
	"github.com/timuraipov/diploma/internal/storage/db"
	"github.com/timuraipov/diploma/internal/usecase"
	"github.com/timuraipov/diploma/pkg/logging"
	"go.uber.org/zap"
)

type OrderWorker struct {
	l            *logging.ZapLogger
	OrderUseCase domain.OrderUseCase
	ctx          context.Context
	cancel       context.CancelFunc
}

func RunWorker(l *logging.ZapLogger, db db.DB, cfg *bootstrap.Config, timeout time.Duration) {
	accrualClient := client.NewClient(cfg.AccrualAddress)
	br := repository.NewBalanceRepository(db)
	balanceUseCase := usecase.NewBalanceUseCase(l, br, timeout)
	orderRepository := repository.NewOrderRepository(db)
	orderUseCase := usecase.NewOrderUseCase(l, orderRepository, accrualClient, balanceUseCase, timeout)
	worker := NewOrderWorker(l, orderUseCase)
	worker.Start()
}

func NewOrderWorker(l *logging.ZapLogger, orderUseCase domain.OrderUseCase) *OrderWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &OrderWorker{
		l:            l,
		OrderUseCase: orderUseCase,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (o *OrderWorker) Start() {
	const numJobs = 5
	jobs := make(chan domain.Order, numJobs)

	go o.StartTicker(o.ctx, o.OrderUseCase, jobs)
	for w := 1; w <= 3; w++ {
		go o.Worker(o.ctx, jobs)
	}
	o.l.InfoCtx(o.ctx, "Starting order worker")
}

func (o *OrderWorker) StartTicker(ctx context.Context, orderUseCase domain.OrderUseCase, jobs chan domain.Order) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			orders, err := orderUseCase.GetUnhandledOrders(ctx)
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
	for {
		select {
		case <-ctx.Done():
			return
		case order := <-jobs:
			statusCode, err := o.OrderUseCase.Accrual(ctx, order)
			if err != nil {
				o.l.ErrorCtx(ctx, "Error processing order- "+err.Error())
			}
			if statusCode == http.StatusTooManyRequests {
				time.Sleep(60 * time.Second)
			}
			continue
		}
	}
}
