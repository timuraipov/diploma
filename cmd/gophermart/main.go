package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	if err := db.RunMigrations(app.Cfg.DSN); err != nil {
		panic(err)
	}
	timeout := time.Duration(app.Cfg.ContextTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// Listen for syscall signals for process to interrupt/quit

	db, err := db.NewDB(ctx, app.Cfg.DSN)
	if err != nil {
		panic(err)
	}
	go bootstrap.RunWorker(l, *db, app.Cfg, timeout)
	r := chi.NewRouter()
	route.Setup(l, app.Cfg, timeout, *db, r)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	srv := &http.Server{
		Addr:    app.Cfg.RunAddress, // app.Cfg.RunAddress
		Handler: r,
	}
	go func() {
		l.InfoCtx(ctx, "Try to start server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.FatalCtx(ctx, "Server error: %s", zap.Error(err))
		}
	}()

	<-stop
	l.InfoCtx(ctx, "Shutting down server...")

	// Контекст с таймаутом
	defer cancel()

	// Корректное завершение
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %s", err)
	}

	l.InfoCtx(ctx, "Server gracefully stopped")
	//err = http.ListenAndServe(app.Cfg.RunAddress, r)

}
