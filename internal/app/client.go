package app

import (
	"context"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/model"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	endpoint       = "http://127.0.0.1:8080/update"
)

var rtMemStats runtime.MemStats

type Client struct {
	pollTicker   *time.Ticker
	reportTicker *time.Ticker
	metrics      map[string]model.Metric
}

func NewClient() *Client {
	m := make(map[string]model.Metric)
	m["Alloc"] = &model.Gauge{}
	m["BuckHashSys"] = &model.Gauge{}
	m["Frees"] = &model.Gauge{}
	m["GCCPUFraction"] = &model.Gauge{}
	m["GCSys"] = &model.Gauge{}
	m["HeapAlloc"] = &model.Gauge{}
	m["HeapIdle"] = &model.Gauge{}
	m["HeapInuse"] = &model.Gauge{}
	m["HeapObjects"] = &model.Gauge{}
	m["HeapReleased"] = &model.Gauge{}
	m["HeapSys"] = &model.Gauge{}
	m["LastGC"] = &model.Gauge{}
	m["Lookups"] = &model.Gauge{}
	m["MCacheInuse"] = &model.Gauge{}
	m["MCacheSys"] = &model.Gauge{}
	m["MSpanInuse"] = &model.Gauge{}
	m["MSpanSys"] = &model.Gauge{}
	m["Mallocs"] = &model.Gauge{}
	m["NextGC"] = &model.Gauge{}
	m["NumForcedGC"] = &model.Gauge{}
	m["NumGC"] = &model.Gauge{}
	m["OtherSys"] = &model.Gauge{}
	m["PauseTotalNs"] = &model.Gauge{}
	m["StackInuse"] = &model.Gauge{}
	m["StackSys"] = &model.Gauge{}
	m["Sys"] = &model.Gauge{}
	m["TotalAlloc"] = &model.Gauge{}

	m["RandomValue"] = &model.Gauge{}
	m["PollCount"] = &model.Counter{}

	return &Client{metrics: m,
		pollTicker:   time.NewTicker(pollInterval),
		reportTicker: time.NewTicker(reportInterval),
	}
}

func (c *Client) updateMetrics() {
	// Заполнение переменной rtMemStats данными статистики
	runtime.ReadMemStats(&rtMemStats)

	//Заполняем поля структуры данными статистики
	c.metrics["Alloc"].SetValue(rtMemStats.Alloc)
	c.metrics["BuckHashSys"].SetValue(rtMemStats.BuckHashSys)
	c.metrics["Frees"].SetValue(rtMemStats.Frees)
	c.metrics["GCCPUFraction"].SetValue(rtMemStats.GCCPUFraction)
	c.metrics["GCSys"].SetValue(rtMemStats.GCSys)
	c.metrics["HeapAlloc"].SetValue(rtMemStats.HeapAlloc)
	c.metrics["HeapIdle"].SetValue(rtMemStats.HeapIdle)
	c.metrics["HeapInuse"].SetValue(rtMemStats.HeapInuse)
	c.metrics["HeapObjects"].SetValue(rtMemStats.HeapObjects)
	c.metrics["HeapReleased"].SetValue(rtMemStats.HeapReleased)
	c.metrics["HeapSys"].SetValue(rtMemStats.HeapSys)
	c.metrics["LastGC"].SetValue(rtMemStats.LastGC)
	c.metrics["Lookups"].SetValue(rtMemStats.Lookups)
	c.metrics["MCacheInuse"].SetValue(rtMemStats.MCacheInuse)
	c.metrics["MCacheSys"].SetValue(rtMemStats.MCacheSys)
	c.metrics["MSpanInuse"].SetValue(rtMemStats.MSpanInuse)
	c.metrics["MSpanSys"].SetValue(rtMemStats.MSpanSys)
	c.metrics["Mallocs"].SetValue(rtMemStats.Mallocs)
	c.metrics["NextGC"].SetValue(rtMemStats.NextGC)
	c.metrics["NumForcedGC"].SetValue(rtMemStats.NumForcedGC)
	c.metrics["NumGC"].SetValue(rtMemStats.NumGC)
	c.metrics["OtherSys"].SetValue(rtMemStats.OtherSys)
	c.metrics["PauseTotalNs"].SetValue(rtMemStats.PauseTotalNs)
	c.metrics["StackInuse"].SetValue(rtMemStats.StackInuse)
	c.metrics["StackSys"].SetValue(rtMemStats.StackSys)
	c.metrics["Sys"].SetValue(rtMemStats.Sys)
	c.metrics["TotalAlloc"].SetValue(rtMemStats.TotalAlloc)

	//Заполняем дополнительные метрики
	count := c.metrics["PollCount"].GetValue().(int64)
	count++
	c.metrics["PollCount"].SetValue(count)
	c.metrics["RandomValue"].SetValue(rand.Float64())
}

func (c *Client) sendMetrics(agent *http.Client) {
	for metricName, metric := range c.metrics {
		url := strings.Join([]string{
			endpoint,
			fmt.Sprintf("%v", metric.GetType()),
			string(metricName),
			fmt.Sprintf("%v", metric.GetValue()),
		},
			"/")
		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			log.Fatal(err)
		}
		request.Header.Add("Content-Type", "text/plain")

		response, err := agent.Do(request)
		if err != nil {
			log.Fatal(err)
		}

		defer response.Body.Close()

		_, err = io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (c *Client) Run(ctx context.Context) {
	transport := &http.Transport{}
	transport.MaxIdleConns = 10

	agent := &http.Client{
		Transport: transport,
		Timeout:   time.Second,
	}

	for {
		select {
		case <-c.pollTicker.C:
			c.updateMetrics()
		case <-c.reportTicker.C:
			c.sendMetrics(agent)
		case <-ctx.Done():
			c.pollTicker.Stop()
			c.reportTicker.Stop()
		}
	}
}
