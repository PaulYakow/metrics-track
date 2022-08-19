package main

import (
	"context"
	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/app/client"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewClientConfig()
	if err != nil {
		log.Println(err)
	}

	client.Run(ctx, cfg)
}
