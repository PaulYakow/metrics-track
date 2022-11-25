package gather

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPSUtil(t *testing.T) {
	valid := []string{"TotalMemory", "FreeMemory", "CPUutilization1"}

	gather := NewGatherPsutil(context.Background())

	namesFromGather := make([]string, 0, len(gather.metrics))
	for name := range gather.metrics {
		namesFromGather = append(namesFromGather, name)
	}

	require.ElementsMatch(t, valid, namesFromGather)
}

func TestNewRuntime(t *testing.T) {
	valid := []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse",
		"HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys",
		"Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc", "RandomValue", "PollCount"}

	gather := NewGatherRuntime()

	namesFromGather := make([]string, 0, len(gather.metrics))
	for name := range gather.metrics {
		namesFromGather = append(namesFromGather, name)
	}

	require.ElementsMatch(t, valid, namesFromGather)
}
