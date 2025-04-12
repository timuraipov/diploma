package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/pkg/logging"
	"go.uber.org/zap"
)

type OrderController struct {
	l            *logging.ZapLogger
	orderUseCase domain.OrderUseCase
	cfg          *bootstrap.Config
}

func NewOrderController(l *logging.ZapLogger, orderUseCase domain.OrderUseCase, cfg *bootstrap.Config) *OrderController {
	return &OrderController{
		l:            l,
		orderUseCase: orderUseCase,
		cfg:          cfg,
	}
}

func (o *OrderController) Save(w http.ResponseWriter, r *http.Request) {
	userId := int64(2)

	orderId, err := io.ReadAll(bufio.NewReader(r.Body))
	if err != nil {
		o.l.ErrorCtx(r.Context(), "invalid data")
	}
	order := domain.Order{
		Number:     string(orderId),
		Status:     "NEW",
		Accrual:    0,
		UserId:     userId,
		UploadedAt: time.Now().Format(time.RFC3339),
	}
	fmt.Print(order)
	w.WriteHeader(http.StatusOK)
}
func (o *OrderController) GetOrders(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("x-user-id").(int64)
	o.l.InfoCtx(r.Context(), "userId", zap.Any("userId", userId))
	orders, err := o.orderUseCase.GetAll(r.Context(), userId)
	if err != nil {
		o.l.ErrorCtx(r.Context(), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

// func (o *OrderController) GetAll() ([]usecase.Order, error) {
// 	orders, err := o.orderUseCase.GetAll(r.con
// 	if err != nil {
// 		o.l.ErrorCtx(o.cfg.ContextTimeout, "failed to get all orders", logging.Error(err))
// 		return nil, err
// 	}
// 	return orders, nil
// }
