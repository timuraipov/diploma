package bootstrap

import (
	"context"
	"time"

	"github.com/timuraipov/diploma/internal/storage/db"
	"github.com/timuraipov/diploma/internal/worker"
	"github.com/timuraipov/diploma/pkg/logging"
	"go.uber.org/zap"
)

type Application struct {
	Cfg     *Config
	DB      *db.DB
	Ctx     context.Context
	Cancel  context.CancelFunc
	Logger  *logging.ZapLogger
	Timeout time.Duration
}

func App() (Application, error) {
	l, err := logging.NewZapLogger(zap.InfoLevel)
	if err != nil {
		panic(err)
	}
	cfg, err := MustLoad()
	if err != nil {
		return Application{}, err
	}
	timeout := time.Duration(cfg.ContextTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	db, err := db.NewDB(ctx, cfg.DSN)
	if err != nil {
		panic(err)
	}
	go worker.RunWorker(l, *db, cfg.AccrualAddress, timeout)
	app := &Application{
		Logger:  l,
		Cfg:     cfg,
		DB:      db,
		Ctx:     ctx,
		Timeout: timeout,
		Cancel:  cancel,
	}

	return *app, nil
}
