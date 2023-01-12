// Package config содержит структуры с конфигурацией сервера.
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

// Config конфигурация сервера (привязка переменных окружения).
type Config struct {
	Address         string        `env:"ADDRESS"`
	GRPCAddress     string        `env:"GRPC_ADDRESS"`
	StoreFile       string        `env:"STORE_FILE" env-default:"/tmp/devops-metrics-db.json"`
	Key             string        `env:"KEY"`
	Dsn             string        `env:"DATABASE_DSN"`
	PathToCryptoKey string        `env:"CRYPTO_KEY"`
	TrustedSubnet   string        `env:"TRUSTED_SUBNET"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL" env-default:"300s"`
	Restore         bool          `env:"RESTORE" env-default:"true"`
}

// Структуры для обработки значений из флагов
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

var grpcAddress = struct {
	name         string
	shorthand    string
	value        *string
	defaultValue string
}{
	"grpc",
	"g",
	new(string),
	"",
}

var storeInterval = struct {
	value        *time.Duration
	name         string
	shorthand    string
	defaultValue time.Duration
}{
	new(time.Duration),
	"interval",
	"i",
	300 * time.Second,
}

var storeFile = struct {
	name         string
	shorthand    string
	value        *string
	defaultValue string
}{
	"file",
	"f",
	new(string),
	"/tmp/devops-metrics-db.json",
}

var restore = struct {
	value        *bool
	name         string
	shorthand    string
	defaultValue bool
}{
	new(bool),
	"restore",
	"r",
	true,
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

var dsn = struct {
	name         string
	shorthand    string
	value        *string
	defaultValue string
}{
	"dsn",
	"d",
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

var trustedSubnet = struct {
	name         string
	shorthand    string
	value        *string
	defaultValue string
}{
	"trusted",
	"t",
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

// NewServerConfig - создаёт объект Config.
func NewServerConfig() (*Config, error) {
	cfg := &Config{}

	cfg.updateCfgFromFlags()

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	fmt.Println("final config =", *cfg)
	return cfg, nil
}

func (cfg *Config) UseHTTPServer() bool {
	return cfg.Address != ""
}

func (cfg *Config) UseGRPCServer() bool {
	return cfg.GRPCAddress != ""
}

func (cfg *Config) updateCfgFromFlags() {
	address.value = pflag.StringP(address.name, address.shorthand, address.defaultValue, "address of HTTP-server in host:port format")
	grpcAddress.value = pflag.StringP(grpcAddress.name, grpcAddress.shorthand, grpcAddress.defaultValue, "grpc address in :port format")
	storeInterval.value = pflag.DurationP(storeInterval.name, storeInterval.shorthand, storeInterval.defaultValue, "store interval in seconds")
	storeFile.value = pflag.StringP(storeFile.name, storeFile.shorthand, storeFile.defaultValue, "path to file")
	restore.value = pflag.BoolP(restore.name, restore.shorthand, restore.defaultValue, "restore after restart")
	hashKey.value = pflag.StringP(hashKey.name, hashKey.shorthand, hashKey.defaultValue, "hash key")
	dsn.value = pflag.StringP(dsn.name, dsn.shorthand, dsn.defaultValue, "DSN for database connect")
	pathToCryptoKey.value = pflag.StringP(pathToCryptoKey.name, pathToCryptoKey.shorthand, pathToCryptoKey.defaultValue, "path to crypto key")
	trustedSubnet.value = pflag.StringP(trustedSubnet.name, trustedSubnet.shorthand, trustedSubnet.defaultValue, "trusted subnet (CIDR notation)")

	pathToConfig.value = pflag.StringP(pathToConfig.name, pathToConfig.shorthand, pathToConfig.defaultValue, "path to config file")

	pflag.Parse()

	cfg.updateCfgFromJSON(*pathToConfig.value)

	if *address.value != address.defaultValue || cfg.Address == "" {
		cfg.Address = *address.value
	}

	if *storeInterval.value != storeInterval.defaultValue || cfg.StoreInterval == 0 {
		cfg.StoreInterval = *storeInterval.value
	}

	if *storeFile.value != storeFile.defaultValue || cfg.StoreFile == "" {
		cfg.StoreFile = *storeFile.value
	}

	if *restore.value != restore.defaultValue {
		cfg.Restore = *restore.value
	}

	if *dsn.value != dsn.defaultValue || cfg.Dsn == "" {
		cfg.Dsn = *dsn.value
	}

	if *pathToCryptoKey.value != pathToCryptoKey.defaultValue || cfg.PathToCryptoKey == "" {
		cfg.PathToCryptoKey = *pathToCryptoKey.value
	}

	if *trustedSubnet.value != trustedSubnet.defaultValue || cfg.TrustedSubnet == "" {
		cfg.TrustedSubnet = *trustedSubnet.value
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
		cfg.StoreInterval = cfgFromJSON.StoreInterval.Duration
		cfg.StoreFile = cfgFromJSON.StoreFile
		cfg.Restore = cfgFromJSON.Restore
		cfg.Dsn = cfgFromJSON.Dsn
		cfg.PathToCryptoKey = cfgFromJSON.PathToCryptoKey
		cfg.TrustedSubnet = cfgFromJSON.TrustedSubnet
	}

	fmt.Println("JSON config =", cfgFromJSON)
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

// CfgFromJSON промежуточная конфигурация сервера (привязка значений из JSON-файла).
type CfgFromJSON struct {
	Address         string   `json:"address"`
	StoreFile       string   `json:"store_file"`
	Dsn             string   `json:"database_dsn"`
	PathToCryptoKey string   `json:"crypto_key"`
	TrustedSubnet   string   `json:"trusted_subnet"`
	StoreInterval   Duration `json:"store_interval"`
	Restore         bool     `json:"restore"`
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
