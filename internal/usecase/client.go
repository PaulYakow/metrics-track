package usecase

import "github.com/PaulYakow/metrics-track/internal/usecase/services/gather"

type ClientUseCase struct {
	repo   ClientMetricRepo
	gather ClientMetricGather
}

func NewClientUC(r ClientMetricRepo) *ClientUseCase {
	return &ClientUseCase{repo: r, gather: gather.New()}
}

func (m *ClientUseCase) Poll() {
	m.repo.Store(m.gather.Update())
}

func (m *ClientUseCase) Report() []string {
	return m.repo.ReadCurrentMetrics()
}
