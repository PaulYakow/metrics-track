package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/pflag"
)

type ServerCfg struct {
	Address       string        `env:"ADDRESS" env-default:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" env-default:"300s"`
	StoreFile     string        `env:"STORE_FILE" env-default:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" env-default:"true"`
	Key           string        `env:"KEY" env-default:""`
	Dsn           string        `env:"DATABASE_DSN" env-default:""`
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

var serverKey = struct {
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

func (cfg *ServerCfg) updateCfgFromFlags() {
	serverAddress.value = pflag.StringP(serverAddress.name, serverAddress.shorthand, serverAddress.defaultValue, "address of server in host:port format")
	storeInterval.value = pflag.DurationP(storeInterval.name, storeInterval.shorthand, storeInterval.defaultValue, "store interval in seconds")
	storeFile.value = pflag.StringP(storeFile.name, storeFile.shorthand, storeFile.defaultValue, "path to file")
	restore.value = pflag.BoolP(restore.name, restore.shorthand, restore.defaultValue, "restore after restart")
	serverKey.value = pflag.StringP(serverKey.name, serverKey.shorthand, serverKey.defaultValue, "hash key")
	dsn.value = pflag.StringP(dsn.name, dsn.shorthand, dsn.defaultValue, "DSN for database connect")

	pflag.Parse()

	cfg.Address = *serverAddress.value
	cfg.StoreInterval = *storeInterval.value
	cfg.StoreFile = *storeFile.value
	cfg.Restore = *restore.value
	cfg.Key = *serverKey.value
	cfg.Dsn = *dsn.value
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
