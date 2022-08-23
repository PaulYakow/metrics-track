package config

import "github.com/spf13/pflag"

func isFlagPassed(name string) bool {
	found := false
	pflag.Visit(func(f *pflag.Flag) {
		if f.Name == name {
			found = true
			return
		}
	})
	return found
}
