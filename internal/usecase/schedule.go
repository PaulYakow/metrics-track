package usecase

import "github.com/PaulYakow/metrics-track/internal/entity"

// Реализация планировщика

type Schedule struct {
	fileRepo   IServerFile
	memoryRepo IServerMemory
}

func NewScheduleUC(file IServerFile, memory IServerMemory) *Schedule {
	return &Schedule{
		fileRepo:   file,
		memoryRepo: memory,
	}
}

func (s *Schedule) RunStoring() {
	gauges, counters := s.memoryRepo.ReadAll()
	metrics := make([]entity.Metrics, 0)

	for name, gauge := range gauges {
		metrics = append(metrics, entity.Metrics{
			ID:    name,
			MType: "gauge",
			Value: gauge.GetPointer(),
		})
	}

	for name, counter := range counters {
		metrics = append(metrics, entity.Metrics{
			ID:    name,
			MType: "counter",
			Delta: counter.GetPointer(),
		})
	}

	s.fileRepo.SaveMetrics(metrics)
}

func (s *Schedule) InitMetrics() {
	s.memoryRepo.InitializeMetrics(s.fileRepo.ReadMetrics())
}
