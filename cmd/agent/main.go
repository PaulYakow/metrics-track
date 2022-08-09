package main

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := app.NewClient()
	client.Run(ctx)
}
