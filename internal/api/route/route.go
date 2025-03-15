package route

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/storage/db"
)

func Setup(cfg *bootstrap.Config, timeout time.Duration, db db.DB, r *chi.Mux) {
	r.Group(func(r chi.Router) {
		NewLoginRouter(cfg, timeout, db, r)
	})
}
