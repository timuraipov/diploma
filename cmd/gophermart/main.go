package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/api/route"
	"github.com/timuraipov/diploma/internal/storage/db"
)

func main() {

	app, err := bootstrap.App()
	if err != nil {
		panic(err)
	}
	fmt.Print(app.Cfg.ContextTimeout, "timeout")
	timeout := time.Duration(app.Cfg.ContextTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_ = cancel //TODO
	db, err := db.NewDB(ctx, app.Cfg.DSN)
	if err != nil {
		panic(err)
	}
	r := chi.NewRouter()
	route.Setup(app.Cfg, timeout, *db, r)
	http.ListenAndServe(app.Cfg.RunAddress, r)
}
