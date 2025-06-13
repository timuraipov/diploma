package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/api/controller"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/repository"
	"github.com/timuraipov/diploma/internal/storage/db"
	"github.com/timuraipov/diploma/internal/usecase"
	"github.com/timuraipov/diploma/pkg/logging"
)

func NewRegisterRouter(logger *logging.ZapLogger, cfg *bootstrap.Config, db db.DB, router chi.Router) {
	authSecret := domain.AuthSecret{
		AccessTokenSecret:      cfg.AccessTokenSecret,
		RefreshTokenSecret:     cfg.RefreshTokenSecret,
		AccessTokenExpiryHour:  cfg.AccessTokenExpiryHour,
		RefreshTokenExpiryHour: cfg.RefreshTokenExpiryHour,
	}
	ur := repository.NewUserRepository(db)
	registerUsecase := usecase.NewRegisterUsecase(logger, ur, authSecret)
	lc := controller.NewRegisterController(logger, registerUsecase, cfg)
	router.Post("/api/user/register", lc.Register)
}
