package usecase

import "github.com/PaulYakow/metrics-track/internal/usecase/services/gather"

type MetricUseCase struct {
	repo   MetricRepo
	gather MetricGather
}

func New(r MetricRepo) *MetricUseCase {
	return &MetricUseCase{repo: r, gather: gather.New()}
}

func (m MetricUseCase) Poll() {
	m.repo.Store(m.gather.Update())
}

func (m MetricUseCase) Report() []string {
	return m.repo.GetCurrentMetrics()
}
