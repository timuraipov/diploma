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

func NewRegisterRouter(logger *logging.ZapLogger, cfg *bootstrap.Config, timeout time.Duration, db db.DB, router chi.Router) {
	lr := repository.NewUserRepository(db)
	registerUsecase := usecase.NewRegisterUsecase(logger, lr, cfg)
	lc := controller.NewRegisterController(logger, registerUsecase, cfg)
	router.Post("/api/user/register", lc.Register)
}
