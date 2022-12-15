// Package config содержит структуры с конфигурациями клиента и сервера.
package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/pflag"
)

// ClientCfg конфигурация клиента.
type ClientCfg struct {
	Address         string        `env:"ADDRESS" env-default:"localhost:8080"`
	Key             string        `env:"KEY" env-default:""`
	ReportInterval  time.Duration `env:"REPORT_INTERVAL" env-default:"10s"`
	PollInterval    time.Duration `env:"POLL_INTERVAL" env-default:"2s"`
	PathToCryptoKey string        `env:"CRYPTO_KEY" env-default:""`
}

var clientAddress = struct {
	name         string
	shorthand    string
	value        *string
	defaultValue string
}{
	"address",
	"a",
	new(string),
	"localhost:8080",
}

var reportInterval = struct {
	value        *time.Duration
	name         string
	shorthand    string
	defaultValue time.Duration
}{
	new(time.Duration),
	"report",
	"r",
	10 * time.Second,
}

var pollInterval = struct {
	value        *time.Duration
	name         string
	shorthand    string
	defaultValue time.Duration
}{
	new(time.Duration),
	"poll",
	"p",
	2 * time.Second,
}

var clientKey = struct {
	name         string
	shorthand    string
	value        *string
	defaultValue string
}{
	"key",
	"k",
	new(string),
	"",
}

var pathToCryptoKey = struct {
	name         string
	shorthand    string
	value        *string
	defaultValue string
}{
	"crypto-key",
	"c",
	new(string),
	"",
}

func (cfg *ClientCfg) updateCfgFromFlags() {
	clientAddress.value = pflag.StringP(clientAddress.name, clientAddress.shorthand, clientAddress.defaultValue, "address of client in host:port format")
	reportInterval.value = pflag.DurationP(reportInterval.name, reportInterval.shorthand, reportInterval.defaultValue, "report interval in seconds")
	pollInterval.value = pflag.DurationP(pollInterval.name, pollInterval.shorthand, pollInterval.defaultValue, "poll interval in seconds")
	clientKey.value = pflag.StringP(clientKey.name, clientKey.shorthand, clientKey.defaultValue, "hash key")
	pathToCryptoKey.value = pflag.StringP(pathToCryptoKey.name, pathToCryptoKey.shorthand, pathToCryptoKey.defaultValue, "path to crypto key")

	pflag.Parse()

	cfg.Address = *clientAddress.value
	cfg.ReportInterval = *reportInterval.value
	cfg.PollInterval = *pollInterval.value
	cfg.Key = *clientKey.value
	cfg.PathToCryptoKey = *pathToCryptoKey.value
}

// NewClientConfig - создаёт объект ClientCfg.
func NewClientConfig() (*ClientCfg, error) {
	cfg := &ClientCfg{}

	cfg.updateCfgFromFlags()

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
