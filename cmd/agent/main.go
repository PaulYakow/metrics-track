package main

import (
	"log"

	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/app/client"
)

func main() {
	cfg, err := config.NewClientConfig()
	if err != nil {
		log.Printf("agent - create config: %v", err)
		return
	}

	client.Run(cfg)
}
