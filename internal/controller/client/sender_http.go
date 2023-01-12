package client

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/PaulYakow/metrics-track/cmd/agent/config"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/httpclient"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/utils/pki"
)

// HTTPSender управляет периодической отправкой метрик на заданный адрес.
type HTTPSender struct {
	client   *httpclient.Client
	uc       usecase.IClient
	logger   logger.ILogger
	endpoint string
	encoder  *pki.Cryptographer
}

// NewHTTPSender создаёт объект HTTPSender.
func NewHTTPSender(client *httpclient.Client, uc usecase.IClient, endpoint string, l logger.ILogger, cfg *config.Config) *HTTPSender {
	s := &HTTPSender{
		client:   client,
		uc:       uc,
		endpoint: endpoint,
		logger:   l,
	}

	if cfg.PathToCryptoKey != "" {
		var err error
		s.encoder, err = pki.NewCryptographer(cfg.PathToCryptoKey)
		if err != nil {
			s.logger.Fatal(err)
		}
	}

	return s
}

// Run - запускает периодическую отправку.
func (s *HTTPSender) Run(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config) {
	ticker := time.NewTicker(cfg.ReportInterval)
	defer wg.Done()

	s.logger.Info("HTTP sender - run with params: address=%s | report=%v | crypto=%s | ip=%s",
		s.endpoint, cfg.ReportInterval, cfg.PathToCryptoKey, cfg.RealIP)
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
func (s *HTTPSender) sendMetricsByURL(routes []string) {
	for _, route := range routes {
		if err := s.client.PostByURL(s.endpoint + route); err != nil {
			s.logger.Error(fmt.Errorf("sender - post metric by URL to %q: %w", s.endpoint+route, err))
		}
	}
}

// sendMetricsByJSON - отправка метрики посредством JSON.
func (s *HTTPSender) sendMetricsByJSON(metrics []entity.Metric) {
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
func (s *HTTPSender) sendMetricsByJSONBatch(metrics []entity.Metric) {
	data, err := json.Marshal(metrics)
	if err != nil {
		s.logger.Error(fmt.Errorf("sender - read metrics: %w", err))
	}

	if s.needEncrypt() {
		if data, err = s.encoder.Encrypt(data); err != nil {
			s.logger.Fatal(err)
		}
	}

	if err = s.client.PostByJSONBatch(s.endpoint, data); err != nil {
		s.logger.Error(fmt.Errorf("sender - post batch of metrics by JSON to %q: %w", s.endpoint, err))
	}
	s.logger.Info("sender - send batch of metrics by JSON: ", string(data))
}

func (s *HTTPSender) needEncrypt() bool {
	return s.encoder != nil
}
