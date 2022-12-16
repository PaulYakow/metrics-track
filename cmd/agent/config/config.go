// Package config содержит структуры с конфигурациями клиента и сервера.
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/pflag"
)

type File struct {
	Path string `env:"CONFIG"`
}

// Config конфигурация клиента.
type Config struct {
	Address         string        `env:"ADDRESS"`
	Key             string        `env:"KEY"`
	ReportInterval  time.Duration `env:"REPORT_INTERVAL"`
	PollInterval    time.Duration `env:"POLL_INTERVAL"`
	PathToCryptoKey string        `env:"CRYPTO_KEY"`
}

var address = struct {
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

var hashKey = struct {
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
	"x",
	new(string),
	"",
}

var pathToConfig = struct {
	name         string
	shorthand    string
	value        *string
	defaultValue string
}{
	"config",
	"c",
	new(string),
	"",
}

// NewClientConfig - создаёт объект Config.
func NewClientConfig() (*Config, error) {
	cfg := &Config{}

	cfg.updateCfgFromFlags()

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) updateCfgFromFlags() {
	address.value = pflag.StringP(address.name, address.shorthand, address.defaultValue, "address of client in host:port format")
	reportInterval.value = pflag.DurationP(reportInterval.name, reportInterval.shorthand, reportInterval.defaultValue, "report interval in seconds")
	pollInterval.value = pflag.DurationP(pollInterval.name, pollInterval.shorthand, pollInterval.defaultValue, "poll interval in seconds")
	hashKey.value = pflag.StringP(hashKey.name, hashKey.shorthand, hashKey.defaultValue, "hash key")
	pathToCryptoKey.value = pflag.StringP(pathToCryptoKey.name, pathToCryptoKey.shorthand, pathToCryptoKey.defaultValue, "path to crypto key")

	pathToConfig.value = pflag.StringP(pathToConfig.name, pathToConfig.shorthand, pathToConfig.defaultValue, "path to config file")

	pflag.Parse()

	cfg.updateCfgFromJSON(*pathToConfig.value)

	if *address.value != address.defaultValue || cfg.Address == "" {
		cfg.Address = *address.value
	}

	if *pollInterval.value != pollInterval.defaultValue || cfg.PollInterval == 0 {
		cfg.PollInterval = *pollInterval.value
	}

	if *reportInterval.value != reportInterval.defaultValue || cfg.ReportInterval == 0 {
		cfg.ReportInterval = *reportInterval.value
	}

	if *pathToCryptoKey.value != pathToCryptoKey.defaultValue || cfg.PathToCryptoKey == "" {
		cfg.PathToCryptoKey = *pathToCryptoKey.value
	}

	cfg.Key = *hashKey.value
}

func (cfg *Config) updateCfgFromJSON(path string) {
	cfgFromJSON := CfgFromJSON{}
	if path == "" {
		path, _ = os.LookupEnv("CONFIG")
	}

	if path != "" {
		if err := cfgFromJSON.loadConfigFromJSON(path); err != nil {
			log.Println("cannot load config file:", err)
		}

		cfg.Address = cfgFromJSON.Address
		cfg.PollInterval = cfgFromJSON.PollInterval.Duration
		cfg.ReportInterval = cfgFromJSON.ReportInterval.Duration
		cfg.PathToCryptoKey = cfgFromJSON.PathToCryptoKey
	}
}
