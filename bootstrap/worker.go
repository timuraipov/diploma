package bootstrap

import (
	"time"

	"github.com/timuraipov/diploma/internal/client"
	"github.com/timuraipov/diploma/internal/repository"
	"github.com/timuraipov/diploma/internal/storage/db"
	"github.com/timuraipov/diploma/internal/usecase"
	"github.com/timuraipov/diploma/internal/worker"
	"github.com/timuraipov/diploma/pkg/logging"
)

func WorkerMustRun(l *logging.ZapLogger, db db.DB, cfg *Config, timeout time.Duration) {
	accrualClient := client.NewClient(cfg.AccrualAddress)
	br := repository.NewBalanceRepository(db)
	balanceUseCase := usecase.NewBalanceUseCase(l, br, timeout)
	orderRepository := repository.NewOrderRepository(db)
	orderUseCase := usecase.NewOrderUseCase(l, orderRepository, *accrualClient, balanceUseCase, timeout)
	worker := worker.NewOrderWorker(l, orderUseCase)
	worker.Start()
}
