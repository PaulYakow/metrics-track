package gather

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
	"math/rand"
	"runtime"
)

var stats runtime.MemStats

type MetricGather struct {
	metrics map[string]*entity.Metric
}

func New() *MetricGather {
	return &MetricGather{
		metrics: initMetrics(),
	}
}

func (g MetricGather) Update() map[string]*entity.Metric {
	runtime.ReadMemStats(&stats)

	g.metrics["Alloc"].UpdateValue(stats.Alloc)
	g.metrics["BuckHashSys"].UpdateValue(stats.BuckHashSys)
	g.metrics["Frees"].UpdateValue(stats.Frees)
	g.metrics["GCCPUFraction"].UpdateValue(stats.GCCPUFraction)
	g.metrics["GCSys"].UpdateValue(stats.GCSys)
	g.metrics["HeapAlloc"].UpdateValue(stats.HeapAlloc)
	g.metrics["HeapIdle"].UpdateValue(stats.HeapIdle)
	g.metrics["HeapInuse"].UpdateValue(stats.HeapInuse)
	g.metrics["HeapObjects"].UpdateValue(stats.HeapObjects)
	g.metrics["HeapReleased"].UpdateValue(stats.HeapReleased)
	g.metrics["HeapSys"].UpdateValue(stats.HeapSys)
	g.metrics["LastGC"].UpdateValue(stats.LastGC)
	g.metrics["Lookups"].UpdateValue(stats.Lookups)
	g.metrics["MCacheInuse"].UpdateValue(stats.MCacheInuse)
	g.metrics["MCacheSys"].UpdateValue(stats.MCacheSys)
	g.metrics["MSpanInuse"].UpdateValue(stats.MSpanInuse)
	g.metrics["MSpanSys"].UpdateValue(stats.MSpanSys)
	g.metrics["Mallocs"].UpdateValue(stats.Mallocs)
	g.metrics["NextGC"].UpdateValue(stats.NextGC)
	g.metrics["NumForcedGC"].UpdateValue(stats.NumForcedGC)
	g.metrics["NumGC"].UpdateValue(stats.NumGC)
	g.metrics["OtherSys"].UpdateValue(stats.OtherSys)
	g.metrics["PauseTotalNs"].UpdateValue(stats.PauseTotalNs)
	g.metrics["StackInuse"].UpdateValue(stats.StackInuse)
	g.metrics["StackSys"].UpdateValue(stats.StackSys)
	g.metrics["Sys"].UpdateValue(stats.Sys)
	g.metrics["TotalAlloc"].UpdateValue(stats.TotalAlloc)
	g.metrics["RandomValue"].UpdateValue(rand.Float64())

	g.metrics["PollCount"].UpdateDelta(1)

	return g.metrics
}

func initMetrics() map[string]*entity.Metric {
	result := make(map[string]*entity.Metric)

	result["Alloc"], _ = entity.Create("gauge", "Alloc", "0")
	result["BuckHashSys"], _ = entity.Create("gauge", "Alloc", "0")
	result["Frees"], _ = entity.Create("gauge", "Frees", "0")
	result["GCCPUFraction"], _ = entity.Create("gauge", "GCCPUFraction", "0")
	result["GCSys"], _ = entity.Create("gauge", "GCSys", "0")
	result["HeapAlloc"], _ = entity.Create("gauge", "HeapAlloc", "0")
	result["HeapIdle"], _ = entity.Create("gauge", "HeapIdle", "0")
	result["HeapInuse"], _ = entity.Create("gauge", "HeapInuse", "0")
	result["HeapObjects"], _ = entity.Create("gauge", "HeapObjects", "0")
	result["HeapReleased"], _ = entity.Create("gauge", "HeapReleased", "0")
	result["HeapSys"], _ = entity.Create("gauge", "HeapSys", "0")
	result["LastGC"], _ = entity.Create("gauge", "LastGC", "0")
	result["Lookups"], _ = entity.Create("gauge", "Lookups", "0")
	result["MCacheInuse"], _ = entity.Create("gauge", "MCacheInuse", "0")
	result["MCacheSys"], _ = entity.Create("gauge", "MCacheSys", "0")
	result["MSpanInuse"], _ = entity.Create("gauge", "MSpanInuse", "0")
	result["MSpanSys"], _ = entity.Create("gauge", "MSpanSys", "0")
	result["Mallocs"], _ = entity.Create("gauge", "Mallocs", "0")
	result["NextGC"], _ = entity.Create("gauge", "NextGC", "0")
	result["NumForcedGC"], _ = entity.Create("gauge", "NumForcedGC", "0")
	result["NumGC"], _ = entity.Create("gauge", "NumGC", "0")
	result["OtherSys"], _ = entity.Create("gauge", "OtherSys", "0")
	result["PauseTotalNs"], _ = entity.Create("gauge", "PauseTotalNs", "0")
	result["StackInuse"], _ = entity.Create("gauge", "StackInuse", "0")
	result["StackSys"], _ = entity.Create("gauge", "StackSys", "0")
	result["Sys"], _ = entity.Create("gauge", "Sys", "0")
	result["TotalAlloc"], _ = entity.Create("gauge", "TotalAlloc", "0")
	result["RandomValue"], _ = entity.Create("gauge", "RandomValue", "0")

	result["PollCount"], _ = entity.Create("counter", "PollCount", "0")

	return result
}
