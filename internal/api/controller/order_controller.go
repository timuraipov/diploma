package controller

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/render"
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
	userID := r.Context().Value(domain.UserIDHeader).(int64)
	orderID, err := io.ReadAll(bufio.NewReader(r.Body))

	if err != nil {
		o.l.ErrorCtx(r.Context(), "invalid data")
	}
	o.l.InfoCtx(r.Context(), "token given", zap.Any("order---", orderID))
	if string(orderID) == "" {
		o.l.ErrorCtx(r.Context(), "orderId is empty")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	o.l.InfoCtx(r.Context(), "orderId", zap.Any("orderId", orderID))
	order := domain.Order{
		ID:         string(orderID),
		Status:     domain.REGISTERED,
		Accrual:    0,
		UserID:     userID,
		UploadedAt: time.Now(), //time.Now().Format(time.RFC3339),
	}

	err = o.orderUseCase.Save(r.Context(), order)
	o.l.InfoCtx(r.Context(), "trying to save order", zap.Any("order", order), zap.Error(err))
	if err != nil {
		if errors.Is(err, domain.ErrOrderAlreadyInProcessing) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, domain.ErrOrderAlreadyProcessedByAnotherUser) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		if errors.Is(err, domain.ErrIncorrectOrderIDFormat) {
			o.l.ErrorCtx(r.Context(), "incorrect orderID format")
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}
	w.WriteHeader(http.StatusAccepted)
}
func (o *OrderController) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(domain.UserIDHeader).(int64)
	o.l.InfoCtx(r.Context(), "userID", zap.Int64("userID", userID))
	orders, err := o.orderUseCase.GetAll(r.Context(), userID)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		o.l.ErrorCtx(r.Context(), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		render.NoContent(w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, orders)

}
