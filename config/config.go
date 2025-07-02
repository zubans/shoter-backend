package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"shuter-go/pkg/logger"
	"strings"
)

type Config struct {
	RunAddr        string `env:"RUN_ADDRESS"`
	DBCfg          string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Migrations     string
}

func NewServerConfig() *Config {
	var db string
	var cfg Config
	var addr string

	flag.StringVar(&addr, "a", "192.168.1.36:8080", "address and port to run server")
	flag.StringVar(&db, "d", "postgres://db_user:db_pass@localhost:5432/mydb?sslmode=disable", "db credential")

	flag.Parse()

	cfg.RunAddr = addr
	cfg.DBCfg = db

	if err := env.Parse(&cfg); err != nil {
		logger.Log.Fatal("Error parsing env vars", zap.Error(err))
	}

	cfg.RunAddr = strings.TrimPrefix(cfg.RunAddr, "http://")
	cfg.RunAddr = strings.TrimPrefix(cfg.RunAddr, "https://")

	cfg.Migrations = "internal/storage/migrations"

	return &cfg
}

func RecoveryServer() {
	if r := recover(); r != nil {
		logger.Log.Fatal("CRITICAL error", zap.Any("error", r))
	}
}
