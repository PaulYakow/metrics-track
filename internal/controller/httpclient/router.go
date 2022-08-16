package httpclient

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/imroc/req/v3"
	"log"
	"time"
)

const (
	pollTime   = 2 * time.Second
	reportTime = 5 * time.Second
	endpoint   = "http://127.0.0.1:8080/update"
)

type clientRoutes struct {
	uc           usecase.ClientMetric
	pollTicker   *time.Ticker
	reportTicker *time.Ticker
}

func NewRouter(ctx context.Context, client *req.Client, uc usecase.ClientMetric) {
	r := &clientRoutes{
		uc:           uc,
		pollTicker:   time.NewTicker(pollTime),
		reportTicker: time.NewTicker(reportTime),
	}

	for {
		select {
		case <-r.pollTicker.C:
			r.uc.Poll()
		case <-r.reportTicker.C:
			r.sendMetrics(client, r.uc.Report())
		case <-ctx.Done():
			r.pollTicker.Stop()
			r.reportTicker.Stop()
		}
	}
}

func (r *clientRoutes) sendMetrics(client *req.Client, routes []string) {
	for _, route := range routes {
		resp, err := client.R().
			SetHeader("Content-Type", "plain/text").
			Post(endpoint + route)

		if err != nil {
			log.Fatal(err)
		}
		resp.Close = true
	}
}
