package gather

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"time"
)

var percent []float64

type gatherPsutil struct {
	metrics map[string]*entity.Metric
}

func NewGatherPsutil() *gatherPsutil {
	gather := new(gatherPsutil)
	gather.initMetrics()
	go gather.getCPUInfo(context.Background())

	return gather
}

func (g *gatherPsutil) initMetrics() {
	g.metrics = make(map[string]*entity.Metric)
	g.metrics["TotalMemory"], _ = entity.Create("gauge", "TotalMemory", "0")
	g.metrics["FreeMemory"], _ = entity.Create("gauge", "FreeMemory", "0")
	g.metrics["CPUutilization1"], _ = entity.Create("gauge", "CPUutilization1", "0")
}

func (g *gatherPsutil) Update() map[string]*entity.Metric {
	memInfo, _ := mem.VirtualMemory()

	g.metrics["TotalMemory"].UpdateValue(memInfo.Total)
	g.metrics["FreeMemory"].UpdateValue(memInfo.Free)
	if len(percent) != 0 {
		g.metrics["CPUutilization1"].UpdateValue(percent[0])
	}

	return g.metrics
}

func (g *gatherPsutil) getCPUInfo(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			percent, _ = cpu.Percent(10*time.Second, false)
		}
	}
}
