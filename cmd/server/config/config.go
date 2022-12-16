package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/pflag"
)

// Config конфигурация сервера.
type Config struct {
	Address         string        `env:"ADDRESS" env-default:"localhost:8080"`
	StoreFile       string        `env:"STORE_FILE" env-default:"/tmp/devops-metrics-db.json"`
	Key             string        `env:"KEY" env-default:""`
	Dsn             string        `env:"DATABASE_DSN" env-default:""`
	PathToCryptoKey string        `env:"CRYPTO_KEY" env-default:""`
	StoreInterval   time.Duration `env:"STORE_INTERVAL" env-default:"300s"`
	Restore         bool          `env:"RESTORE" env-default:"true"`
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

func (cfg *Config) updateCfgFromFlags() {
	address.value = pflag.StringP(address.name, address.shorthand, address.defaultValue, "address of server in host:port format")
	storeInterval.value = pflag.DurationP(storeInterval.name, storeInterval.shorthand, storeInterval.defaultValue, "store interval in seconds")
	storeFile.value = pflag.StringP(storeFile.name, storeFile.shorthand, storeFile.defaultValue, "path to file")
	restore.value = pflag.BoolP(restore.name, restore.shorthand, restore.defaultValue, "restore after restart")
	hashKey.value = pflag.StringP(hashKey.name, hashKey.shorthand, hashKey.defaultValue, "hash key")
	dsn.value = pflag.StringP(dsn.name, dsn.shorthand, dsn.defaultValue, "DSN for database connect")
	pathToCryptoKey.value = pflag.StringP(pathToCryptoKey.name, pathToCryptoKey.shorthand, pathToCryptoKey.defaultValue, "path to crypto key")

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

	cfg.Key = *hashKey.value
	fmt.Println("flag config =", *cfg)
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
	}

	fmt.Println("JSON config =", cfgFromJSON)
}
