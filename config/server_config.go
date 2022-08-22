package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type ServerCfg struct {
	Address       string        `env:"ADDRESS" env-default:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" env-default:"300s"`
	StoreFile     string        `env:"STORE_FILE" env-default:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" env-default:"true"`
}

func NewServerConfig() (*ServerCfg, error) {
	cfg := &ServerCfg{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
