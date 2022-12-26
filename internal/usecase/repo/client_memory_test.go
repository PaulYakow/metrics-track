package repo

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PaulYakow/metrics-track/internal/entity"
)

func TestClientMemoryRepo(t *testing.T) {
	var repo *ClientMemoryRepo
	var metrics map[string]*entity.Metric

	t.Run("client repo create", func(t *testing.T) {
		repo = NewClientRepo()
		require.IsType(t, &ClientMemoryRepo{}, repo)
	})

	t.Run("client repo store new", func(t *testing.T) {
		var delta int64 = 999
		var value float64 = 0.999
		metrics = map[string]*entity.Metric{
			"testGauge": {
				ID:    "testGauge",
				MType: "gauge",
				Value: &value},
			"testCounter": {
				ID:    "testCounter",
				MType: "counter",
				Delta: &delta},
		}
		repo.Store(metrics)
		require.Equal(t, metrics, repo.metrics)
	})

	t.Run("client repo store update", func(t *testing.T) {
		var value float64 = 0.666
		metricsUpdate := map[string]*entity.Metric{
			"testGauge": {
				ID:    "testGauge",
				MType: "gauge",
				Value: &value},
		}
		repo.Store(metricsUpdate)
		require.Equal(t, repo.metrics["testGauge"], metricsUpdate["testGauge"])
	})

	t.Run("client repo read all", func(t *testing.T) {
		allMetrics := repo.ReadAll()
		require.NotEmpty(t, allMetrics)
		for _, metric := range allMetrics {
			require.Equal(t, *metrics[metric.ID], metric)
		}
	})
}
