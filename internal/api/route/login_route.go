package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/api/controller"
	"github.com/timuraipov/diploma/internal/repository"
	"github.com/timuraipov/diploma/internal/storage/db"
	"github.com/timuraipov/diploma/internal/usecase"
	"github.com/timuraipov/diploma/pkg/logging"
)

func NewLoginRouter(logger *logging.ZapLogger, cfg *bootstrap.Config, db db.DB, router chi.Router) {
	lr := repository.NewUserRepository(db)
	lu := usecase.NewLoginUsecase(logger, lr, cfg)
	lc := controller.NewLoginController(logger, lu, cfg)
	router.Post("/api/user/login", lc.Login)
}
