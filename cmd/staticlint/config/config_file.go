//go:build config

package config

import (
	"embed"
	"encoding/json"
	"log"
)

//go:embed config.json
var config embed.FS

func NewCfgData() ConfigData {
	data, err := config.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
