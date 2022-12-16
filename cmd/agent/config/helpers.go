package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Duration struct {
	time.Duration
}

func (duration *Duration) UnmarshalJSON(b []byte) error {
	var unmarshalledJson interface{}

	err := json.Unmarshal(b, &unmarshalledJson)
	if err != nil {
		return err
	}

	switch value := unmarshalledJson.(type) {
	case float64:
		duration.Duration = time.Duration(value)
	case string:
		duration.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid duration: %#v", unmarshalledJson)
	}

	return nil
}

type CfgFromJSON struct {
	Address         string   `json:"address"`
	ReportInterval  Duration `json:"report_interval"`
	PollInterval    Duration `json:"poll_interval"`
	PathToCryptoKey string   `json:"crypto_key"`
}

func (cfg *CfgFromJSON) loadConfigFromJSON(path string) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}
	err = json.NewDecoder(file).Decode(cfg)
	return err
}
