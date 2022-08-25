package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/pflag"
	"time"
)

type ServerCfg struct {
	Address       string        `env:"ADDRESS" env-default:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" env-default:"300s"`
	StoreFile     string        `env:"STORE_FILE" env-default:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" env-default:"true"`
}

var serverAddress = struct {
	name         string
	shorthand    string
	value        *string
	defaultValue string
}{
	"serverAddress",
	"a",
	new(string),
	"localhost:8080",
}

var storeInterval = struct {
	name         string
	shorthand    string
	value        *time.Duration
	defaultValue time.Duration
}{
	"interval",
	"i",
	new(time.Duration),
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
	name         string
	shorthand    string
	value        *bool
	defaultValue bool
}{
	"restore",
	"r",
	new(bool),
	true,
}

func (cfg *ServerCfg) updateCfgFromFlags() {
	serverAddress.value = pflag.StringP(serverAddress.name, serverAddress.shorthand, serverAddress.defaultValue, "address of server in host:port format")
	storeInterval.value = pflag.DurationP(storeInterval.name, storeInterval.shorthand, storeInterval.defaultValue, "store interval in seconds")
	storeFile.value = pflag.StringP(storeFile.name, storeFile.shorthand, storeFile.defaultValue, "path to file")
	restore.value = pflag.BoolP(restore.name, restore.shorthand, restore.defaultValue, "restore after restart")

	pflag.Parse()

	if isFlagPassed(serverAddress.name) {
		cfg.Address = *serverAddress.value
	}

	if isFlagPassed(storeInterval.name) {
		cfg.StoreInterval = *storeInterval.value
	}

	if isFlagPassed(storeFile.name) {
		cfg.StoreFile = *storeFile.value
	}

	if isFlagPassed(restore.name) {
		cfg.Restore = *restore.value
	}
}

func NewServerConfig() (*ServerCfg, error) {
	cfg := &ServerCfg{}

	cfg.updateCfgFromFlags()

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
