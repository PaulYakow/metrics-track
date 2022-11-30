// Package client содержит логику периодического сбора и отправки метрик на заданный адрес.
package client

import (
	"context"
	"sync"
	"time"

	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
)

// Collector управляет периодическим сбором метрик.
type Collector struct {
	uc     usecase.IClient
	logger logger.ILogger
}

// NewCollector создаёт объект Collector.
func NewCollector(uc usecase.IClient, l logger.ILogger) *Collector {
	return &Collector{
		uc:     uc,
		logger: l,
	}
}

// Run - запускает периодический сбор метрик.
func (c *Collector) Run(ctx context.Context, wg *sync.WaitGroup, interval time.Duration) {
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
