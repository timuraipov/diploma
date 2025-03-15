package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddress     string `env:"RUN_ADDRESS"`
	DatabaseURI    string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

// -a-d-r
func MustLoad() (*Config, error) {
	cfg := &Config{}
	flag.StringVar(&cfg.RunAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.DatabaseURI, "d", "postgres://metric:XXXXX@localhost:5432/metric?sslmode=disable", "database dsn")
	flag.StringVar(&cfg.AccrualAddress, "r", "", "accrual system address")

	flag.Parse()
	err := env.Parse(cfg)
	return cfg, err
}
