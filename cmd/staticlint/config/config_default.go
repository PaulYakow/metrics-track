//go:build !config

package config

func NewCfgData() ConfigData {
	return ConfigData{
		Staticcheck: []string{
			"S1001",
			"ST1000",
			"QF1003",
		},
	}
}
