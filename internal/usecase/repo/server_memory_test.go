package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PaulYakow/metrics-track/internal/entity"
)

func TestServerMemoryRepo(t *testing.T) {
	var repo *ServerMemoryRepo
	var metric *entity.Metric
	var metrics []entity.Metric

	t.Run("server repo create", func(t *testing.T) {
		repo = NewServerMemory()
		require.IsType(t, &ServerMemoryRepo{}, repo)
	})

	t.Run("server repo store new", func(t *testing.T) {
		var value float64 = 0.999
		metric = &entity.Metric{
			ID:    "testGauge",
			MType: "gauge",
			Value: &value}
		err := repo.Store(metric)
		require.NoError(t, err)
		require.Equal(t, repo.metrics["testGauge"], metric)
	})

	t.Run("server repo store unknown type", func(t *testing.T) {
		metricUnknown := &entity.Metric{
			ID:    "testUnknown",
			MType: "unknown"}
		err := repo.Store(metricUnknown)
		require.Error(t, err)
		require.Empty(t, repo.metrics[metricUnknown.ID])
	})

	t.Run("server repo store batch", func(t *testing.T) {
		var value float64 = 0.999
		var delta int64 = 999
		repo = NewServerMemory()

		metrics = []entity.Metric{
			{ID: "testGauge", MType: "gauge", Value: &value},
			{ID: "testCounter", MType: "counter", Delta: &delta},
		}
		err := repo.StoreBatch(metrics)
		require.NoError(t, err)
		require.NotEmpty(t, repo.metrics)
		for _, metric := range metrics {
			require.Equal(t, *repo.metrics[metric.ID], metric)
		}
	})

	t.Run("server repo read single", func(t *testing.T) {
		exist, err := repo.Read(context.Background(), entity.Metric{ID: "testGauge"})
		require.NoError(t, err)
		require.NotEmpty(t, exist)
	})

	t.Run("server repo read unknown", func(t *testing.T) {
		exist, err := repo.Read(context.Background(), entity.Metric{ID: "testUnknown"})
		require.Error(t, err)
		require.Empty(t, exist)
	})

	t.Run("server repo read cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		exist, err := repo.Read(ctx, entity.Metric{ID: "testUnknown"})
		require.NoError(t, err)
		require.Empty(t, exist)
	})

	t.Run("server repo read all", func(t *testing.T) {
		allMetrics, err := repo.ReadAll(context.Background())
		require.NoError(t, err)
		for _, metric := range allMetrics {
			require.Equal(t, *repo.metrics[metric.ID], metric)
		}
	})

	t.Run("server repo read all cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		allMetrics, err := repo.ReadAll(ctx)
		require.NoError(t, err)
		require.Empty(t, allMetrics)
	})

	t.Run("server repo connection error", func(t *testing.T) {
		err := repo.CheckConnection()
		require.Error(t, err)
	})
}
