// Package config содержит структуры с конфигурацией клиента.
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/pflag"
)

// Config конфигурация клиента (привязка переменных окружения).
type Config struct {
	Address         string        `env:"ADDRESS" env-default:"localhost:8080"`
	Key             string        `env:"KEY" env-default:""`
	ReportInterval  time.Duration `env:"REPORT_INTERVAL" env-default:"10s"`
	PollInterval    time.Duration `env:"POLL_INTERVAL" env-default:"2s"`
	PathToCryptoKey string        `env:"CRYPTO_KEY"`
	RealIP          string        `env:"REAL_IP" env-default:"127.0.0.1"`
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

var realIP = struct {
	name         string
	shorthand    string
	value        *string
	defaultValue string
}{
	"ip",
	"i",
	new(string),
	"127.0.0.1",
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
	realIP.value = pflag.StringP(realIP.name, realIP.shorthand, realIP.defaultValue, "real IP for header X-Real-IP")

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

	if *realIP.value != realIP.defaultValue || cfg.RealIP == "" {
		cfg.RealIP = *realIP.value
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
		cfg.RealIP = cfgFromJSON.RealIP
	}
}

type Duration struct {
	time.Duration
}

func (duration *Duration) UnmarshalJSON(b []byte) error {
	var unmarshalled interface{}

	err := json.Unmarshal(b, &unmarshalled)
	if err != nil {
		return err
	}

	switch value := unmarshalled.(type) {
	case float64:
		duration.Duration = time.Duration(value)
	case string:
		duration.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid duration: %#v", unmarshalled)
	}

	return nil
}

// CfgFromJSON промежуточная конфигурация клиента (привязка значений из JSON-файла).
type CfgFromJSON struct {
	Address         string   `json:"address"`
	ReportInterval  Duration `json:"report_interval"`
	PollInterval    Duration `json:"poll_interval"`
	PathToCryptoKey string   `json:"crypto_key"`
	RealIP          string   `json:"real_ip"`
}

func (cfg *CfgFromJSON) loadConfigFromJSON(path string) error {
	file, err := os.Open(path)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("close config file: %v\n", err)
		}
	}(file)

	if err != nil {
		return err
	}
	err = json.NewDecoder(file).Decode(cfg)
	return err
}
