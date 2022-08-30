package client

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"sync"
	"time"
)

type collector struct {
	ctx context.Context
	uc  usecase.IClient
	l   logger.ILogger
}

func NewCollector(ctx context.Context, uc usecase.IClient, l logger.ILogger) *collector {
	return &collector{
		ctx: ctx,
		uc:  uc,
		l:   l,
	}
}

func (c *collector) Run(wg *sync.WaitGroup, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer wg.Done()

	c.l.Info("collector - run with interval %v", interval)
	for {
		select {
		case <-ticker.C:
			c.uc.Poll()
		case <-c.ctx.Done():
			ticker.Stop()
			c.l.Info("collector - context %v", c.ctx.Err())
			return
		}
	}
}
