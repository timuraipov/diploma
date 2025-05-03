package controller

import (
	"bufio"
	"encoding/json"
	"errors"
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
	userId := r.Context().Value("x-user-id").(int64)
	orderId, err := io.ReadAll(bufio.NewReader(r.Body))
	if err != nil {
		o.l.ErrorCtx(r.Context(), "invalid data")
	}
	o.l.InfoCtx(r.Context(), "token given", zap.Any("order---", orderId))
	if string(orderId) == "" {
		o.l.ErrorCtx(r.Context(), "orderId is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	o.l.InfoCtx(r.Context(), "orderId", zap.Any("orderId", orderId))
	order := domain.Order{
		ID:         string(orderId),
		Status:     domain.REGISTERED,
		Accrual:    0,
		UserId:     userId,
		UploadedAt: time.Now(), //time.Now().Format(time.RFC3339),
	}

	err = o.orderUseCase.Save(r.Context(), order)
	o.l.InfoCtx(r.Context(), "trying to save order", zap.Any("order", order), zap.Error(err))
	if err != nil {
		if errors.Is(err, domain.OrderAlreadyInProcessing) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, domain.OrderAlreadyProcessedByAnotherUser) {
			w.WriteHeader(http.StatusConflict)
		}
	}
	w.WriteHeader(http.StatusAccepted)
}
func (o *OrderController) GetOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("into GetOrders")
	userId := r.Context().Value("x-user-id").(int64)
	fmt.Println("userid", userId)
	o.l.InfoCtx(r.Context(), "userId", zap.Int64("userId", userId))
	orders, err := o.orderUseCase.GetAll(r.Context(), userId)
	if err != nil {
		o.l.ErrorCtx(r.Context(), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
