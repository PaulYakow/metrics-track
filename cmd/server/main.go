package main

import (
	"log"

	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/app/server"
)

func main() {
	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("server - create config: %v\n", err)
	}

	server.Run(cfg)
}
