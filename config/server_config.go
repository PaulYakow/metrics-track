package config

import "github.com/ilyakaznacheev/cleanenv"

type ServerCfg struct {
	Address []string `env:"ADDRESS" env-separator:":" env-default:"localhost:8080"`
}

func NewServerConfig() (*ServerCfg, error) {
	cfg := &ServerCfg{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
