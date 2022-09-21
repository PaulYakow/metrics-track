package client

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"sync"
	"time"
)

type collector struct {
	uc     usecase.IClient
	logger logger.ILogger
}

func NewCollector(uc usecase.IClient, l logger.ILogger) *collector {
	return &collector{
		uc:     uc,
		logger: l,
	}
}

func (c *collector) Run(ctx context.Context, wg *sync.WaitGroup, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer wg.Done()

	c.logger.Info("collector - run with params: p=%v", interval)
	for {
		select {
		case <-ticker.C:
			c.uc.Poll()
		case <-ctx.Done():
			ticker.Stop()
			c.logger.Info("collector - context %v", ctx.Err())
			return
		}
	}
}
