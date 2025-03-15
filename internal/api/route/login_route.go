package route

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/api/controller"
	"github.com/timuraipov/diploma/internal/repository"
	"github.com/timuraipov/diploma/internal/storage/db"
	"github.com/timuraipov/diploma/internal/usecase"
)

func NewLoginRouter(cfg *bootstrap.Config, timeout time.Duration, db db.DB, router chi.Router) {
	lr := repository.NewLoginRepository(db)
	lc := &controller.LoginController{
		LoginUsecase: usecase.NewLoginUsecase(lr, timeout),
		Cfg:          cfg,
	}
	router.Post("/api/user/login", lc.Login)
}
