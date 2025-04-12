package route

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/api/middleware"
	"github.com/timuraipov/diploma/internal/storage/db"
	"github.com/timuraipov/diploma/pkg/logging"
)

func Setup(logger *logging.ZapLogger, cfg *bootstrap.Config, timeout time.Duration, db db.DB, r *chi.Mux) {
	r.Group(func(r chi.Router) {
		NewLoginRouter(logger, cfg, timeout, db, r)
		NewRegisterRouter(logger, cfg, timeout, db, r)
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware(cfg.AccessTokenSecret))
		NewOrderRouter(logger, cfg, timeout, db, r)
	})
}
