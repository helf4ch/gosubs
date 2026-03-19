package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Addr            string `env:"ADDR,required"`
	DbConn          string `env:"DB_CONN,required"`
	GooseMigrations string `env:"GOOSE_MIGRATION_DIR,required"`
}

func New() (*Config, error) {
	cfg := &Config{}

	err := envconfig.Process(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
