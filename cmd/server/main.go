package main

import (
	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/app/server"
	"log"
)

func main() {
	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Printf("create config for server: %v", err)
	}

	server.Run(cfg)
}
