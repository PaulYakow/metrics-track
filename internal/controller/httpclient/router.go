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

type Client struct {
	pollTicker   *time.Ticker
	reportTicker *time.Ticker
	uc           usecase.Metric
}

func New(ucm usecase.Metric) *Client {
	return &Client{
		pollTicker:   time.NewTicker(pollTime),
		reportTicker: time.NewTicker(reportTime),
		uc:           ucm,
	}
}

func (c *Client) Run(ctx context.Context, client *req.Client) {
	for {
		select {
		case <-c.pollTicker.C:
			c.uc.Poll()
		case <-c.reportTicker.C:
			c.sendMetrics(client, c.uc.Report())
		case <-ctx.Done():
			c.pollTicker.Stop()
			c.reportTicker.Stop()
		}
	}
}

func (c *Client) sendMetrics(client *req.Client, routes []string) {
	for _, route := range routes {
		_, err := client.R().
			SetHeader("Content-Type", "plain/text").
			Post(endpoint + route)
		if err != nil {
			log.Fatal(err)
		}
	}
}
