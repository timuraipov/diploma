package route

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/api/controller"
	"github.com/timuraipov/diploma/internal/client"
	"github.com/timuraipov/diploma/internal/repository"
	"github.com/timuraipov/diploma/internal/storage/db"
	"github.com/timuraipov/diploma/internal/usecase"
	"github.com/timuraipov/diploma/pkg/logging"
)

func NewOrderRouter(logger *logging.ZapLogger, cfg *bootstrap.Config, timeout time.Duration, db db.DB, router chi.Router) {
	accrualClient := client.NewClient(cfg.AccrualAddress)
	br := repository.NewBalanceRepository(db)
	balanceUseCase := usecase.NewBalanceUseCase(logger, br, timeout)
	or := repository.NewOrderRepository(db)
	orderUseCase := usecase.NewOrderUseCase(logger, or, accrualClient, balanceUseCase, timeout)
	OrderController := controller.NewOrderController(logger, orderUseCase, cfg)
	router.Get("/api/user/orders", OrderController.GetOrders)
	router.Post("/api/user/orders", OrderController.Save)
}
