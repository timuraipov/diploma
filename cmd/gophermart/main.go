package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/api/route"
	"github.com/timuraipov/diploma/internal/storage/db"
	"github.com/timuraipov/diploma/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	l, err := logging.NewZapLogger(zap.InfoLevel)
	if err != nil {
		panic(err)
	}
	app, err := bootstrap.App()
	if err != nil {
		panic(err)
	}
	timeout := time.Duration(app.Cfg.ContextTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_ = cancel //TODO
	db, err := db.NewDB(ctx, app.Cfg.DSN)
	if err != nil {
		panic(err)
	}
	r := chi.NewRouter()
	route.Setup(l, app.Cfg, timeout, *db, r)
	l.InfoCtx(ctx, "Try to start server")
	err = http.ListenAndServe(app.Cfg.RunAddress, r)
	if err != nil {
		l.PanicCtx(ctx, "failed to start server", zap.Error(err))
	}
}
