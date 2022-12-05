//go:build !config

package main

func NewCfgData() ConfigData {
	return ConfigData{
		Staticcheck: []string{
			"S1001",
			"ST1000",
			"QF1003",
		},
	}
}
