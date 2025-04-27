package route

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/api/controller"
	"github.com/timuraipov/diploma/internal/repository"
	"github.com/timuraipov/diploma/internal/storage/db"
	"github.com/timuraipov/diploma/internal/usecase"
	"github.com/timuraipov/diploma/pkg/logging"
)

func NewBalanceRouter(logger *logging.ZapLogger, cfg *bootstrap.Config, timeout time.Duration, db db.DB, router chi.Router) {
	balanceRepository := repository.NewBalanceRepository(db)
	balanceUseCase := usecase.NewBalanceUseCase(logger, balanceRepository, timeout)
	balanceController := controller.NewBalanceController(logger, balanceUseCase, cfg)
	router.Get("/api/user/balance", balanceController.GetBalance)
	router.Get("/api/user/withdrawals", balanceController.Withdrawals)
	router.Post("/api/user/balance/withdraw", balanceController.Withdraw)
}
