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
	userID := r.Context().Value(domain.UserIDHeader).(int64)
	balance, err := b.balanceUseCase.GetBalance(r.Context(), userID)
	if err != nil {
		b.l.ErrorCtx(r.Context(), err.Error())
		render.Status(r, http.StatusBadGateway)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, balance)
}

func (b *BalanceController) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(domain.UserIDHeader).(int64)
	var request domain.WithdrawRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}
	withdraw := domain.Withdraw{
		ID:          request.Order,
		Sum:         request.Sum,
		UserID:      userID,
		ProcessedAt: time.Now(), //time.Now().Format(time.RFC3339),
	}
	err = b.balanceUseCase.Withdraw(r.Context(), withdraw)
	if err != nil {
		if errors.Is(err, domain.ErrInsufficientFunds) {
			render.Status(r, http.StatusPaymentRequired)
			return
		}
		if errors.Is(err, domain.ErrWithdrawAlreadyUsed) {
			render.Status(r, http.StatusUnprocessableEntity)
			return
		} //TODO add error handle
		b.l.ErrorCtx(r.Context(), err.Error())
		render.Status(r, http.StatusBadGateway)
		//	w.WriteHeader(http.StatusBadGateway)
		return
	}
	render.Status(r, http.StatusOK)
}

func (b *BalanceController) Withdrawals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(domain.UserIDHeader).(int64)
	withdrawals, err := b.balanceUseCase.Withdrawals(r.Context(), userID)

	if err != nil {
		b.l.InfoCtx(r.Context(), err.Error())
		render.Status(r, http.StatusBadGateway)
	}
	if len(withdrawals) == 0 {
		w.Header().Set("Content-Type", "application/json") // не знаю зачем тесты падают.
		render.Status(r, http.StatusNoContent)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, withdrawals)
}
