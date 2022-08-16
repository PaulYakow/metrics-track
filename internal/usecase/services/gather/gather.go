package gather

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
	"math/rand"
	"runtime"
)

var stats runtime.MemStats

type MetricGather struct {
	gauges   map[string]*entity.Gauge
	counters map[string]*entity.Counter
}

func New() *MetricGather {
	return &MetricGather{
		gauges:   initGauges(),
		counters: initCounters(),
	}
}

func (g MetricGather) Update() (map[string]*entity.Gauge, map[string]*entity.Counter) {
	g.updateGauges()
	g.updateCounters()
	return g.gauges, g.counters
}

func initGauges() map[string]*entity.Gauge {
	result := make(map[string]*entity.Gauge)

	result["Alloc"] = &entity.Gauge{}
	result["BuckHashSys"] = &entity.Gauge{}
	result["Frees"] = &entity.Gauge{}
	result["GCCPUFraction"] = &entity.Gauge{}
	result["GCSys"] = &entity.Gauge{}
	result["HeapAlloc"] = &entity.Gauge{}
	result["HeapIdle"] = &entity.Gauge{}
	result["HeapInuse"] = &entity.Gauge{}
	result["HeapObjects"] = &entity.Gauge{}
	result["HeapReleased"] = &entity.Gauge{}
	result["HeapSys"] = &entity.Gauge{}
	result["LastGC"] = &entity.Gauge{}
	result["Lookups"] = &entity.Gauge{}
	result["MCacheInuse"] = &entity.Gauge{}
	result["MCacheSys"] = &entity.Gauge{}
	result["MSpanInuse"] = &entity.Gauge{}
	result["MSpanSys"] = &entity.Gauge{}
	result["Mallocs"] = &entity.Gauge{}
	result["NextGC"] = &entity.Gauge{}
	result["NumForcedGC"] = &entity.Gauge{}
	result["NumGC"] = &entity.Gauge{}
	result["OtherSys"] = &entity.Gauge{}
	result["PauseTotalNs"] = &entity.Gauge{}
	result["StackInuse"] = &entity.Gauge{}
	result["StackSys"] = &entity.Gauge{}
	result["Sys"] = &entity.Gauge{}
	result["TotalAlloc"] = &entity.Gauge{}
	result["RandomValue"] = &entity.Gauge{}

	return result
}

func initCounters() map[string]*entity.Counter {
	result := make(map[string]*entity.Counter)

	result["PollCount"] = &entity.Counter{}

	return result
}

func (g MetricGather) updateGauges() {
	runtime.ReadMemStats(&stats)

	g.gauges["Alloc"].SetValue(stats.Alloc)
	g.gauges["BuckHashSys"].SetValue(stats.BuckHashSys)
	g.gauges["Frees"].SetValue(stats.Frees)
	g.gauges["GCCPUFraction"].SetValue(stats.GCCPUFraction)
	g.gauges["GCSys"].SetValue(stats.GCSys)
	g.gauges["HeapAlloc"].SetValue(stats.HeapAlloc)
	g.gauges["HeapIdle"].SetValue(stats.HeapIdle)
	g.gauges["HeapInuse"].SetValue(stats.HeapInuse)
	g.gauges["HeapObjects"].SetValue(stats.HeapObjects)
	g.gauges["HeapReleased"].SetValue(stats.HeapReleased)
	g.gauges["HeapSys"].SetValue(stats.HeapSys)
	g.gauges["LastGC"].SetValue(stats.LastGC)
	g.gauges["Lookups"].SetValue(stats.Lookups)
	g.gauges["MCacheInuse"].SetValue(stats.MCacheInuse)
	g.gauges["MCacheSys"].SetValue(stats.MCacheSys)
	g.gauges["MSpanInuse"].SetValue(stats.MSpanInuse)
	g.gauges["MSpanSys"].SetValue(stats.MSpanSys)
	g.gauges["Mallocs"].SetValue(stats.Mallocs)
	g.gauges["NextGC"].SetValue(stats.NextGC)
	g.gauges["NumForcedGC"].SetValue(stats.NumForcedGC)
	g.gauges["NumGC"].SetValue(stats.NumGC)
	g.gauges["OtherSys"].SetValue(stats.OtherSys)
	g.gauges["PauseTotalNs"].SetValue(stats.PauseTotalNs)
	g.gauges["StackInuse"].SetValue(stats.StackInuse)
	g.gauges["StackSys"].SetValue(stats.StackSys)
	g.gauges["Sys"].SetValue(stats.Sys)
	g.gauges["TotalAlloc"].SetValue(stats.TotalAlloc)
	g.gauges["RandomValue"].SetValue(rand.Float64())
}

func (g MetricGather) updateCounters() {
	g.counters["PollCount"].Increment()
}
