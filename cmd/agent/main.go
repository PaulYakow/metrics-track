package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type gauge float64
type counter int64

var ms runtime.MemStats

var pollInterval = 2 * time.Second
var reportInterval = 10 * time.Second
var endpoint = "http://127.0.0.1:8080/update/"

type rtMetrics struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	PollCount     counter
	RandomValue   gauge
}

func (rm *rtMetrics) update() {
	//Заполняем переменную ms данными статистики
	runtime.ReadMemStats(&ms)

	//Заполняем поля структуры данными статистики
	rm.Alloc = gauge(ms.Alloc)
	rm.BuckHashSys = gauge(ms.BuckHashSys)
	rm.Frees = gauge(ms.Frees)
	rm.GCCPUFraction = gauge(ms.GCCPUFraction)
	rm.GCSys = gauge(ms.GCSys)
	rm.HeapAlloc = gauge(ms.HeapAlloc)
	rm.HeapIdle = gauge(ms.HeapIdle)
	rm.HeapInuse = gauge(ms.HeapInuse)
	rm.HeapObjects = gauge(ms.HeapObjects)
	rm.HeapReleased = gauge(ms.HeapReleased)
	rm.HeapSys = gauge(ms.HeapSys)
	rm.LastGC = gauge(ms.LastGC)
	rm.Lookups = gauge(ms.Lookups)
	rm.MCacheInuse = gauge(ms.MCacheInuse)
	rm.MCacheSys = gauge(ms.MCacheSys)
	rm.MSpanInuse = gauge(ms.MSpanInuse)
	rm.MSpanSys = gauge(ms.MSpanSys)
	rm.Mallocs = gauge(ms.Mallocs)
	rm.NextGC = gauge(ms.NextGC)
	rm.NumForcedGC = gauge(ms.NumForcedGC)
	rm.NumGC = gauge(ms.NumGC)
	rm.OtherSys = gauge(ms.OtherSys)
	rm.PauseTotalNs = gauge(ms.PauseTotalNs)
	rm.StackInuse = gauge(ms.StackInuse)
	rm.StackSys = gauge(ms.StackSys)
	rm.Sys = gauge(ms.Sys)
	rm.TotalAlloc = gauge(ms.TotalAlloc)

	//Заполняем дополнительные метрики
	rm.PollCount++
	rm.RandomValue = gauge(rand.Float64())

	time.Sleep(pollInterval)
}

func (rm *rtMetrics) post(agent *http.Client) {
	request, err := http.NewRequest(http.MethodPost, endpoint+"gauge/Alloc/"+fmt.Sprintf("%v", rm.Alloc), nil)
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Add("Content-Type", "text/plain")

	response, err := agent.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	_, err = io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(reportInterval)
}

func main() {
	var metrics rtMetrics

	transport := &http.Transport{}
	transport.MaxIdleConns = 10

	agent := &http.Client{
		Transport: transport,
		Timeout:   time.Second,
	}

	for {
		go metrics.update()
		go metrics.post(agent)
	}
}
