package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

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
		w.WriteHeader(http.StatusBadGateway)
		return // TODO add different statuses
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(balance) //todo encode Handle
	if err != nil {
		b.l.ErrorCtx(r.Context(), err.Error())
	}
}

func (b *BalanceController) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(domain.UserIDHeader).(int64)
	var request domain.WithdrawRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
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
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}
		if errors.Is(err, domain.ErrWithdrawAlreadyUsed) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		} //TODO add error handle
		b.l.ErrorCtx(r.Context(), err.Error())
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (b *BalanceController) Withdrawals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(domain.UserIDHeader).(int64)
	withdrawals, err := b.balanceUseCase.Withdrawals(r.Context(), userID)

	if err != nil {
		b.l.InfoCtx(r.Context(), err.Error())
		w.WriteHeader(http.StatusBadGateway)
	}
	if len(withdrawals) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(withdrawals)
	if err != nil {
		b.l.ErrorCtx(r.Context(), err.Error())
	}
}

// func (b *BalanceController) GetBalance(w http.ResponseWriter, r *http.Request) {
// 	userID := r.Context().Value("x-user-id").(int64)
// 	b.l.InfoCtx(r.Context(), "GetBalance", zap.Int64("userID", userID))
// 	balance, err := b.balanceUseCase.GetBalance(r.Context(), userID)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
// 	balanceResponse := &domain.BalanceResponse{
// 		Current:   balance.Balance,
// 		Withdrawn: 0,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(balanceResponse)
// }
// func (b *BalanceController) Withdraw(w http.ResponseWriter, r *http.Request) {
// 	userID := r.Context().Value("x-user-id").(int64)
// 	b.l.InfoCtx(r.Context(), "GetBalance", zap.Int64("userID", userID))
// 	var request domain.WithdrawRequest
// 	err := json.NewDecoder(r.Body).Decode(&request)
// 	if err != nil {
// 		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
// 		return
// 	}
// 	withdraw := domain.Withdraw{
// 		ID:          request.Order,
// 		Sum:         request.Sum,
// 		UserId:      userID,
// 		ProcessedAt: time.Now(),
// 	}
// 	err = b.balanceUseCase.Withdraw(r.Context(), withdraw)

// }
