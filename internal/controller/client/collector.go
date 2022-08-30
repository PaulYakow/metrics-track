package client

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"log"
	"time"
)

type collector struct {
	ctx context.Context
	uc  usecase.IClient
}

func NewCollector(ctx context.Context, uc usecase.IClient) *collector {
	return &collector{
		ctx: ctx,
		uc:  uc,
	}
}

func (c *collector) Run(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("collector run with interval %v", interval)
	for {
		select {
		case <-ticker.C:
			c.uc.Poll()
			log.Printf("polling... %v", time.Now())
		case <-c.ctx.Done():
			ticker.Stop()
			log.Printf("collector context %v", c.ctx.Err())
			return
		}
	}
}
