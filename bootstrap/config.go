package bootstrap

import (
	"flag"
	"sync"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddress             string `env:"RUN_ADDRESS"`
	DSN                    string `env:"DATABASE_URI"`
	AccrualAddress         string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	ContextTimeout         int64  `env:"CONTEXT_TIMEOUT" envDefault:"3"`
	AccessTokenExpiryHour  int    `env:"ACCESS_TOKEN_EXPIRY_HOUR" envDefault:"2"`
	RefreshTokenExpiryHour int    `env:"REFRESH_TOKEN_EXPIRY_HOUR" envDefault:"48"`
	AccessTokenSecret      string `env:"ACCESS_TOKEN_SECRET" envDefault:"secret"`
	RefreshTokenSecret     string `env:"REFRESH_TOKEN_SECRET" envDefault:"secret"`
}

var once sync.Once

func MustLoad() (*Config, error) {
	cfg := &Config{}

	once.Do(func() {
		flag.StringVar(&cfg.RunAddress, "a", "localhost:8081", "address and port to run server")
		flag.StringVar(&cfg.DSN, "d", "postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable", "database dsn")
		flag.StringVar(&cfg.AccrualAddress, "r", "http://localhost:8080", "accrual system address")
		flag.Int64Var(&cfg.ContextTimeout, "t", 3, "time for timeout")
		flag.Parse()
	})

	err := env.Parse(cfg)
	if err != nil {
		panic(err)
	}
	return cfg, nil
}
