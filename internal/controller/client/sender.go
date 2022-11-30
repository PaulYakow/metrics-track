package client

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/httpclient"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
)

// Sender управляет периодической отправкой метрик на заданный адрес.
type Sender struct {
	client   *httpclient.Client
	uc       usecase.IClient
	logger   logger.ILogger
	endpoint string
}

// NewSender создаёт объект Sender.
func NewSender(client *httpclient.Client, uc usecase.IClient, endpoint string, l logger.ILogger) *Sender {
	return &Sender{
		client:   client,
		uc:       uc,
		endpoint: endpoint,
		logger:   l,
	}
}

// Run - запускает периодическую отправку.
func (s *Sender) Run(ctx context.Context, wg *sync.WaitGroup, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer wg.Done()

	s.logger.Info("sender - run with params: a=%s | r=%v", s.endpoint, interval)
	for {
		select {
		case <-ticker.C:
			s.sendMetricsByJSONBatch(s.uc.GetAll())
		case <-ctx.Done():
			ticker.Stop()
			s.logger.Info("sender - context canceled")
			return
		}
	}
}

// sendMetricsByURL - отправка метрики посредством URL.
func (s *Sender) sendMetricsByURL(routes []string) {
	for _, route := range routes {
		if err := s.client.PostByURL(s.endpoint + route); err != nil {
			s.logger.Error(fmt.Errorf("sender - post metric by URL to %q: %w", s.endpoint+route, err))
		}
	}
}

// sendMetricsByJSON - отправка метрики посредством JSON.
func (s *Sender) sendMetricsByJSON(metrics []entity.Metric) {
	for _, metric := range metrics {
		data, err := json.Marshal(metric)
		if err != nil {
			s.logger.Error(fmt.Errorf("sender - read metric %q: %w", metric.ID, err))
		}
		if err = s.client.PostByJSON(s.endpoint, data); err != nil {
			s.logger.Error(fmt.Errorf("sender - post metric by JSON to %q: %w", s.endpoint, err))
		}
	}
}

// sendMetricsByJSONBatch - отправка метрики посредством пакета JSON.
func (s *Sender) sendMetricsByJSONBatch(metrics []entity.Metric) {
	data, err := json.Marshal(metrics)
	if err != nil {
		s.logger.Error(fmt.Errorf("sender - read metrics: %w", err))
	}

	if err = s.client.PostByJSONBatch(s.endpoint, data); err != nil {
		s.logger.Error(fmt.Errorf("sender - post batch of metrics by JSON to %q: %w", s.endpoint, err))
	}
}
