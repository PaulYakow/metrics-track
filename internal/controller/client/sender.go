package client

import (
	"fmt"
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
	l        logger.ILogger
}

func NewSender(client *httpclient.Client, uc usecase.IClient, endpoint string, l logger.ILogger) *sender {
	return &sender{
		client:   client,
		uc:       uc,
		endpoint: endpoint,
		l:        l,
	}
}

func (s *sender) Run(wg *sync.WaitGroup, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer wg.Done()

	s.l.Info("sender - run with interval %v", interval)
	for {
		select {
		case <-ticker.C:
			s.sendMetricsByJSON(s.uc.UpdateValues())
		case <-s.client.Done():
			s.l.Info("sender - context canceled")
			ticker.Stop()
			return
		}
	}
}

func (s *sender) sendMetricsByURL(routes []string) {
	for _, route := range routes {
		if err := s.client.PostByURL(s.endpoint + route); err != nil {
			s.l.Error(fmt.Errorf("sender - post metric by URL to %q: %w", s.endpoint+route, err))
		}
	}
}

func (s *sender) sendMetricsByJSON(data [][]byte) {
	for _, rawMetric := range data {
		if err := s.client.PostByJSON(s.endpoint, rawMetric); err != nil {
			s.l.Error(fmt.Errorf("sender - post metric by JSON to %q: %w", s.endpoint, err))
		}
	}
}
