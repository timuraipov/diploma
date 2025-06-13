package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/api/route"
	"github.com/timuraipov/diploma/internal/storage/db"
	"go.uber.org/zap"
)

func main() {

	app, err := bootstrap.App()
	if err != nil {
		panic(err)
	}
	if err := db.RunMigrations(app.Cfg.DSN); err != nil {
		panic(err)
	}
	// Listen for syscall signals for process to interrupt/quit
	defer app.Cancel()

	defer func() {
		app.Logger.InfoCtx(app.Ctx, "Close DB pool connections", zap.Error(err))
		app.DB.Pool.Close()
	}()
	r := chi.NewRouter()
	route.Setup(app.Logger, app.Cfg, app.Timeout, *app.DB, r)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	srv := &http.Server{
		Addr:    app.Cfg.RunAddress, // app.Cfg.RunAddress
		Handler: r,
	}
	go func() {
		app.Logger.InfoCtx(app.Ctx, "Try to start server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.FatalCtx(app.Ctx, "Server error: %s", zap.Error(err))
		}
	}()

	<-stop
	app.Logger.InfoCtx(app.Ctx, "Shutting down server...")

	// Контекст с таймаутом

	// Корректное завершение
	if err := srv.Shutdown(app.Ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %s", err)
	}

	app.Logger.InfoCtx(app.Ctx, "Server gracefully stopped")
}
