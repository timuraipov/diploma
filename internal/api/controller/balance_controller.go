package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/pkg/logging"
)

type BalanceController struct {
	l              *logging.ZapLogger
	balanceUseCase domain.BalanceUseCase
	cfg            *bootstrap.Config
}

func NewBalanceController(logger *logging.ZapLogger, br domain.BalanceUseCase, cfg *bootstrap.Config) *BalanceController {
	return &BalanceController{
		l:              logger,
		balanceUseCase: br,
		cfg:            cfg,
	}
}

func (b *BalanceController) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(domain.UserIDHeader).(int64)
	if !ok {
		handleErrorResponse(b.l, w, r, domain.ErrUserIDNotFound, http.StatusUnauthorized)
	}
	balance, err := b.balanceUseCase.GetBalance(r.Context(), userID)
	if err != nil {
		handleErrorResponse(b.l, w, r, err, http.StatusBadGateway)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, balance)
}

func (b *BalanceController) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(domain.UserIDHeader).(int64)
	if !ok {
		handleErrorResponse(b.l, w, r, domain.ErrUserIDNotFound, http.StatusUnauthorized)
	}
	var request domain.WithdrawRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handleErrorResponse(b.l, w, r, err, http.StatusBadRequest)
		return
	}
	withdraw := domain.Withdraw{
		OrderID:     request.Order,
		Sum:         request.Sum,
		UserID:      userID,
		ProcessedAt: time.Now(), // time.Now().Format(time.RFC3339),
	}
	err = b.balanceUseCase.Withdraw(r.Context(), withdraw)
	if err != nil {
		if errors.Is(err, domain.ErrInsufficientFunds) {
			handleErrorResponse(b.l, w, r, err, http.StatusPaymentRequired)
			return
		}
		if errors.Is(err, domain.ErrWithdrawAlreadyUsed) {
			handleErrorResponse(b.l, w, r, err, http.StatusUnprocessableEntity)
			return
		}
		handleErrorResponse(b.l, w, r, err, http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (b *BalanceController) Withdrawals(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(domain.UserIDHeader).(int64)
	if !ok {
		handleErrorResponse(b.l, w, r, domain.ErrUserIDNotFound, http.StatusUnauthorized)
	}
	withdrawals, err := b.balanceUseCase.Withdrawals(r.Context(), userID)
	if err != nil {
		b.l.InfoCtx(r.Context(), err.Error())
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if len(withdrawals) == 0 {
		w.Header().Set("Content-Type", "application/json") // не знаю зачем тесты падают.
		render.NoContent(w, r)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, withdrawals)
}
