package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/httpclient"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"sync"
	"time"
)

type sender struct {
	client   *httpclient.Client
	uc       usecase.IClient
	endpoint string
	logger   logger.ILogger
}

func NewSender(client *httpclient.Client, uc usecase.IClient, endpoint string, l logger.ILogger) *sender {
	return &sender{
		client:   client,
		uc:       uc,
		endpoint: endpoint,
		logger:   l,
	}
}

func (s *sender) Run(ctx context.Context, wg *sync.WaitGroup, interval time.Duration) {
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

func (s *sender) sendMetricsByURL(routes []string) {
	for _, route := range routes {
		if err := s.client.PostByURL(s.endpoint + route); err != nil {
			s.logger.Error(fmt.Errorf("sender - post metric by URL to %q: %w", s.endpoint+route, err))
		}
	}
}

func (s *sender) sendMetricsByJSON(metrics []entity.Metric) {
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

func (s *sender) sendMetricsByJSONBatch(metrics []entity.Metric) {
	data, err := json.Marshal(metrics)
	if err != nil {
		s.logger.Error(fmt.Errorf("sender - read metrics: %w", err))
	}

	if err = s.client.PostByJSONBatch(s.endpoint, data); err != nil {
		s.logger.Error(fmt.Errorf("sender - post batch of metrics by JSON to %q: %w", s.endpoint, err))
	}
}
