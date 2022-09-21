package gather

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
	"math/rand"
	"runtime"
)

var stats runtime.MemStats

type gatherRuntime struct {
	metrics map[string]*entity.Metric
}

func NewGatherRuntime() *gatherRuntime {
	gather := new(gatherRuntime)
	gather.initMetrics()
	return gather
}

func (g *gatherRuntime) initMetrics() map[string]*entity.Metric {
	g.metrics = make(map[string]*entity.Metric)

	g.metrics["Alloc"], _ = entity.Create("gauge", "Alloc", "0")
	g.metrics["BuckHashSys"], _ = entity.Create("gauge", "BuckHashSys", "0")
	g.metrics["Frees"], _ = entity.Create("gauge", "Frees", "0")
	g.metrics["GCCPUFraction"], _ = entity.Create("gauge", "GCCPUFraction", "0")
	g.metrics["GCSys"], _ = entity.Create("gauge", "GCSys", "0")
	g.metrics["HeapAlloc"], _ = entity.Create("gauge", "HeapAlloc", "0")
	g.metrics["HeapIdle"], _ = entity.Create("gauge", "HeapIdle", "0")
	g.metrics["HeapInuse"], _ = entity.Create("gauge", "HeapInuse", "0")
	g.metrics["HeapObjects"], _ = entity.Create("gauge", "HeapObjects", "0")
	g.metrics["HeapReleased"], _ = entity.Create("gauge", "HeapReleased", "0")
	g.metrics["HeapSys"], _ = entity.Create("gauge", "HeapSys", "0")
	g.metrics["LastGC"], _ = entity.Create("gauge", "LastGC", "0")
	g.metrics["Lookups"], _ = entity.Create("gauge", "Lookups", "0")
	g.metrics["MCacheInuse"], _ = entity.Create("gauge", "MCacheInuse", "0")
	g.metrics["MCacheSys"], _ = entity.Create("gauge", "MCacheSys", "0")
	g.metrics["MSpanInuse"], _ = entity.Create("gauge", "MSpanInuse", "0")
	g.metrics["MSpanSys"], _ = entity.Create("gauge", "MSpanSys", "0")
	g.metrics["Mallocs"], _ = entity.Create("gauge", "Mallocs", "0")
	g.metrics["NextGC"], _ = entity.Create("gauge", "NextGC", "0")
	g.metrics["NumForcedGC"], _ = entity.Create("gauge", "NumForcedGC", "0")
	g.metrics["NumGC"], _ = entity.Create("gauge", "NumGC", "0")
	g.metrics["OtherSys"], _ = entity.Create("gauge", "OtherSys", "0")
	g.metrics["PauseTotalNs"], _ = entity.Create("gauge", "PauseTotalNs", "0")
	g.metrics["StackInuse"], _ = entity.Create("gauge", "StackInuse", "0")
	g.metrics["StackSys"], _ = entity.Create("gauge", "StackSys", "0")
	g.metrics["Sys"], _ = entity.Create("gauge", "Sys", "0")
	g.metrics["TotalAlloc"], _ = entity.Create("gauge", "TotalAlloc", "0")
	g.metrics["RandomValue"], _ = entity.Create("gauge", "RandomValue", "0")

	g.metrics["PollCount"], _ = entity.Create("counter", "PollCount", "0")

	return g.metrics
}

func (g *gatherRuntime) Update() map[string]*entity.Metric {
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
