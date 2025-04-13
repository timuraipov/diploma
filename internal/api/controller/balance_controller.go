package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/pkg/logging"
	"go.uber.org/zap"
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
	userID := r.Context().Value("x-user-id").(int64)
	b.l.InfoCtx(r.Context(), "GetBalance", zap.Int64("userID", userID))
	balance, err := b.balanceUseCase.GetBalance(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	balanceResponse := &domain.BalanceResponse{
		Current:   balance.Balance,
		Withdrawn: 0,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balanceResponse)
}
func (b *BalanceController) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("x-user-id").(int64)
	b.l.InfoCtx(r.Context(), "GetBalance", zap.Int64("userID", userID))
	var request domain.WithdrawRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}
	withdraw := domain.Withdraw{
		Order:       request.Order,
		Sum:         request.Sum,
		UserId:      userID,
		ProcessedAt: time.Now(),
	}
	err = b.balanceUseCase.Withdraw(r.Context(), withdraw)
}
