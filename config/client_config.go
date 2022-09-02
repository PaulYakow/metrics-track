package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/pflag"
	"time"
)

type ClientCfg struct {
	Address        string        `env:"ADDRESS" env-default:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" env-default:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" env-default:"2s"`
	Key            string        `env:"KEY" env-default:""`
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
	name         string
	shorthand    string
	value        *time.Duration
	defaultValue time.Duration
}{
	"report",
	"r",
	new(time.Duration),
	10 * time.Second,
}

var pollInterval = struct {
	name         string
	shorthand    string
	value        *time.Duration
	defaultValue time.Duration
}{
	"poll",
	"p",
	new(time.Duration),
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

func (cfg *ClientCfg) updateCfgFromFlags() {
	clientAddress.value = pflag.StringP(clientAddress.name, clientAddress.shorthand, clientAddress.defaultValue, "address of client in host:port format")
	reportInterval.value = pflag.DurationP(reportInterval.name, reportInterval.shorthand, reportInterval.defaultValue, "report interval in seconds")
	pollInterval.value = pflag.DurationP(pollInterval.name, pollInterval.shorthand, pollInterval.defaultValue, "poll interval in seconds")
	clientKey.value = pflag.StringP(clientKey.name, clientKey.shorthand, clientKey.defaultValue, "hash key")

	pflag.Parse()

	if isFlagPassed(clientAddress.name) {
		cfg.Address = *clientAddress.value
	}

	if isFlagPassed(reportInterval.name) {
		cfg.ReportInterval = *reportInterval.value
	}

	if isFlagPassed(pollInterval.name) {
		cfg.PollInterval = *pollInterval.value
	}

	if isFlagPassed(clientKey.name) {
		cfg.Key = *clientKey.value
	}
}

func NewClientConfig() (*ClientCfg, error) {
	cfg := &ClientCfg{}

	cfg.updateCfgFromFlags()

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
