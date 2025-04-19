package worker

import (
	"context"

	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/pkg/logging"
)

type OrderWorker struct {
	l               *logging.ZapLogger
	orderRepository domain.OrderRepository
	ctx             context.Context
	cancel          context.CancelFunc
}

func NewOrderWorker(l *logging.ZapLogger, orderRepository domain.OrderRepository) *OrderWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &OrderWorker{
		l:               l,
		orderRepository: orderRepository,
		ctx:             ctx,
		cancel:          cancel,
	}
}
// func (o *OrderWorker) Start() {
// 	o.l.InfoCtx(o.ctx, "Starting order worker")
// 	go func() {
// 		for {
// 			select {
// 			case <-o.ctx.Done():
// 				o.l.InfoCtx(o.ctx, "Stopping order worker")
// 				return
// 			default:
// 				// Process orders here
// 				orders, err := o.orderRepository.GetAll(o.ctx)
// 				if err != nil {
// 					o.l.ErrorCtx(o.ctx, "Error getting orders- "+err.Error())
// 					continue
// 				}
// 				for _, order := range orders {
// 					// Process each order
// 					o.l.InfoCtx(o.ctx, "Processing order", order)
// 				}
// 			}
// 		}
// 	}()
}
