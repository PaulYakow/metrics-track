package main

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/app/client"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client.Run(ctx)
}
