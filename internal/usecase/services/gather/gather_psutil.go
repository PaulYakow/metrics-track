package gather

import (
	"context"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/PaulYakow/metrics-track/internal/entity"
)

var percent []float64

// GatherPsutil реализация "сборщика" метрик psutil (IClientGather).
type GatherPsutil struct {
	metrics map[string]*entity.Metric
}

// NewGatherPsutil создаёт объект GatherPsutil
func NewGatherPsutil(ctx context.Context) *GatherPsutil {
	g := new(GatherPsutil)
	g.metrics = make(map[string]*entity.Metric, 4)

	g.metrics["TotalMemory"], _ = entity.Create("gauge", "TotalMemory", "0")
	g.metrics["FreeMemory"], _ = entity.Create("gauge", "FreeMemory", "0")
	g.metrics["CPUutilization1"], _ = entity.Create("gauge", "CPUutilization1", "0")
	go g.getCPUInfo(ctx)

	return g
}

func (g *GatherPsutil) Update() map[string]*entity.Metric {
	memInfo, _ := mem.VirtualMemory()

	g.metrics["TotalMemory"].UpdateValue(memInfo.Total)
	g.metrics["FreeMemory"].UpdateValue(memInfo.Free)
	if len(percent) != 0 {
		g.metrics["CPUutilization1"].UpdateValue(percent[0])
	}

	return g.metrics
}

func (g *GatherPsutil) getCPUInfo(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("cpu info - %v", ctx.Err())
			return
		default:
			percent, _ = cpu.Percent(10*time.Second, false)
		}
	}
}
