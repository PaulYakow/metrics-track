package httpclient

import (
	"context"
	"fmt"
	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/imroc/req/v3"
	"log"
	"time"
)

type clientRoutes struct {
	uc           usecase.IClient
	endpoint     string
	pollTicker   *time.Ticker
	reportTicker *time.Ticker
}

func NewRouter(ctx context.Context, cfg *config.ClientCfg, client *req.Client, uc usecase.IClient) {
	r := &clientRoutes{
		uc:           uc,
		endpoint:     fmt.Sprintf("http://%s:%s/update", cfg.Address[0], cfg.Address[1]),
		pollTicker:   time.NewTicker(cfg.PollInterval),
		reportTicker: time.NewTicker(cfg.ReportInterval),
	}

	for {
		select {
		case <-r.pollTicker.C:
			r.uc.Poll()
		case <-r.reportTicker.C:
			//r.sendMetricsByURL(client, r.uc.UpdateRoutes())
			r.sendMetricsByJSON(client, r.uc.UpdateValues())
		case <-ctx.Done():
			r.pollTicker.Stop()
			r.reportTicker.Stop()
		}
	}
}

func (r *clientRoutes) sendMetricsByURL(client *req.Client, routes []string) {
	for _, route := range routes {
		resp, err := client.R().
			SetHeader("Content-Type", "plain/text").
			Post(r.endpoint + route)

		if err != nil {
			log.Println(err)
		}
		resp.Close = true
	}
}

func (r *clientRoutes) sendMetricsByJSON(client *req.Client, data [][]byte) {
	request := client.R()

	for _, rawMetric := range data {
		r.sendMetric(request, rawMetric)
	}
}

func (r *clientRoutes) sendMetric(request *req.Request, rawMetric []byte) {
	resp, err := request.
		SetHeader("Content-Type", "application/json").
		SetBody(rawMetric).
		Post(r.endpoint)

	if err != nil {
		fmt.Println("", err)
		log.Println(err)
	}

	defer resp.Body.Close()
}
