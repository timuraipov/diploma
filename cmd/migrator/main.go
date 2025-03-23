package main

import (
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/storage/db"
)

func main() {
	cfg, err := bootstrap.MustLoad()
	if err != nil {
		panic(err)
	}
	if err := db.RunMigrations(cfg.DSN); err != nil {
		panic(err)
	}
}
