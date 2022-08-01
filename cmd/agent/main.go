package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type gauge float64
type counter int64
type name string

var ms runtime.MemStats

var pollInterval = 2 * time.Second
var reportInterval = 10 * time.Second
var endpoint = "http://127.0.0.1:8080/update"

type metric struct {
	tMetric string
	value   any
}

func createMetricBatch() map[name]*metric {
	metrics := make(map[name]*metric)

	metrics["Alloc"] = &metric{
		tMetric: "gauge",
	}
	metrics["BuckHashSys"] = &metric{
		tMetric: "gauge",
	}
	metrics["Frees"] = &metric{
		tMetric: "gauge",
	}
	metrics["GCCPUFraction"] = &metric{
		tMetric: "gauge",
	}
	metrics["GCSys"] = &metric{
		tMetric: "gauge",
	}
	metrics["HeapAlloc"] = &metric{
		tMetric: "gauge",
	}
	metrics["HeapIdle"] = &metric{
		tMetric: "gauge",
	}
	metrics["HeapInuse"] = &metric{
		tMetric: "gauge",
	}
	metrics["HeapObjects"] = &metric{
		tMetric: "gauge",
	}
	metrics["HeapReleased"] = &metric{
		tMetric: "gauge",
	}
	metrics["HeapSys"] = &metric{
		tMetric: "gauge",
	}
	metrics["LastGC"] = &metric{
		tMetric: "gauge",
	}
	metrics["Lookups"] = &metric{
		tMetric: "gauge",
	}
	metrics["MCacheInuse"] = &metric{
		tMetric: "gauge",
	}
	metrics["MCacheSys"] = &metric{
		tMetric: "gauge",
	}
	metrics["MSpanInuse"] = &metric{
		tMetric: "gauge",
	}
	metrics["MSpanSys"] = &metric{
		tMetric: "gauge",
	}
	metrics["Mallocs"] = &metric{
		tMetric: "gauge",
	}
	metrics["NextGC"] = &metric{
		tMetric: "gauge",
	}
	metrics["NumForcedGC"] = &metric{
		tMetric: "gauge",
	}
	metrics["NumGC"] = &metric{
		tMetric: "gauge",
	}
	metrics["OtherSys"] = &metric{
		tMetric: "gauge",
	}
	metrics["PauseTotalNs"] = &metric{
		tMetric: "gauge",
	}
	metrics["StackSys"] = &metric{
		tMetric: "gauge",
	}
	metrics["StackInuse"] = &metric{
		tMetric: "gauge",
	}
	metrics["Sys"] = &metric{
		tMetric: "gauge",
	}
	metrics["TotalAlloc"] = &metric{
		tMetric: "gauge",
	}
	metrics["PollCount"] = &metric{
		tMetric: "counter",
		value:   counter(-1),
	}
	metrics["RandomValue"] = &metric{
		tMetric: "gauge",
	}

	return metrics
}

func updateMetrics(m map[name]*metric) {
	//Заполняем переменную ms данными статистики
	runtime.ReadMemStats(&ms)

	//Заполняем поля структуры данными статистики
	m["Alloc"].value = gauge(ms.Alloc)
	m["BuckHashSys"].value = gauge(ms.BuckHashSys)
	m["Frees"].value = gauge(ms.Frees)
	m["GCCPUFraction"].value = gauge(ms.GCCPUFraction)
	m["GCSys"].value = gauge(ms.GCSys)
	m["HeapAlloc"].value = gauge(ms.HeapAlloc)
	m["HeapIdle"].value = gauge(ms.HeapIdle)
	m["HeapInuse"].value = gauge(ms.HeapInuse)
	m["HeapObjects"].value = gauge(ms.HeapObjects)
	m["HeapReleased"].value = gauge(ms.HeapReleased)
	m["HeapSys"].value = gauge(ms.HeapSys)
	m["LastGC"].value = gauge(ms.LastGC)
	m["Lookups"].value = gauge(ms.Lookups)
	m["MCacheInuse"].value = gauge(ms.MCacheInuse)
	m["MCacheSys"].value = gauge(ms.MCacheSys)
	m["MSpanInuse"].value = gauge(ms.MSpanInuse)
	m["MSpanSys"].value = gauge(ms.MSpanSys)
	m["Mallocs"].value = gauge(ms.Mallocs)
	m["NextGC"].value = gauge(ms.NextGC)
	m["NumForcedGC"].value = gauge(ms.NumForcedGC)
	m["NumGC"].value = gauge(ms.NumGC)
	m["OtherSys"].value = gauge(ms.OtherSys)
	m["PauseTotalNs"].value = gauge(ms.PauseTotalNs)
	m["StackInuse"].value = gauge(ms.StackInuse)
	m["StackSys"].value = gauge(ms.StackSys)
	m["Sys"].value = gauge(ms.Sys)
	m["TotalAlloc"].value = gauge(ms.TotalAlloc)

	//	//Заполняем дополнительные метрики
	count := m["PollCount"].value.(counter)
	count++
	m["PollCount"].value = count
	m["RandomValue"].value = gauge(rand.Float64())
}

func sendReport(agent *http.Client, m map[name]*metric) {
	for mtrName, mtr := range m {
		urlElems := []string{
			endpoint,
			fmt.Sprintf("%v", mtr.tMetric),
			string(mtrName),
			fmt.Sprintf("%v", mtr.value),
		}
		url := strings.Join(urlElems, "/")
		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			fmt.Println("request| ", err)
		}
		request.Header.Add("Content-Type", "text/plain")

		response, err := agent.Do(request)
		if err != nil {
			fmt.Println("response| ", err)
		}
		defer response.Body.Close()
		_, err = io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("response read all| ", err)
		}
	}
}

func main() {
	metrics := createMetricBatch()

	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	transport := &http.Transport{}
	transport.MaxIdleConns = 10

	agent := &http.Client{
		Transport: transport,
		Timeout:   time.Second,
	}

	for {
		select {
		case <-pollTicker.C:
			updateMetrics(metrics)
		case <-reportTicker.C:
			sendReport(agent, metrics)
		}
	}

}
