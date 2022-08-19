package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type ClientCfg struct {
	Address        []string      `env:"ADDRESS" env-separator:":" env-default:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" env-default:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" env-default:"2s"`
}

func NewClientConfig() (*ClientCfg, error) {
	cfg := &ClientCfg{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
