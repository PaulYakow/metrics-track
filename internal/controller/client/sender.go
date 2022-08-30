package client

import (
	"github.com/PaulYakow/metrics-track/internal/pkg/httpclient"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"log"
	"time"
)

type sender struct {
	client   *httpclient.Client
	uc       usecase.IClient
	endpoint string
	ticker   *time.Ticker
}

func NewSender(client *httpclient.Client, uc usecase.IClient, endpoint string, interval time.Duration) *sender {
	log.Printf("new sender with interval %v, endpoint %v", interval, endpoint)
	return &sender{
		client:   client,
		uc:       uc,
		endpoint: endpoint,
		ticker:   time.NewTicker(interval),
	}
}

func (s *sender) Run() {
	for {
		select {
		case <-s.ticker.C:
			s.sendMetricsByJSON(s.uc.UpdateValues())
			log.Printf("sending... %v", time.Now())
		case <-s.client.Done():
			s.ticker.Stop()
			return
		}
	}
}

func (s *sender) sendMetricsByURL(routes []string) {
	for _, route := range routes {
		if err := s.client.PostByURL(s.endpoint + route); err != nil {
			log.Printf("sender - post metric by URL to %q: %v", s.endpoint+route, err)
		}
	}
}

func (s *sender) sendMetricsByJSON(data [][]byte) {
	for _, rawMetric := range data {
		if err := s.client.PostByJSON(s.endpoint, rawMetric); err != nil {
			log.Printf("sender - post metric by JSON to %q: %v", s.endpoint, err)
		}
	}
}
