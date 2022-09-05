package v1

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"time"
)

type scheduler struct {
	memory usecase.IServerMemory
	repo   usecase.IServerRepo
}

func NewScheduler(memory usecase.IServerMemory, repo usecase.IServerRepo) *scheduler {
	return &scheduler{
		memory: memory,
		repo:   repo,
	}
}

func (s *scheduler) Run(ctx context.Context, restore bool, interval time.Duration) {
	if restore {
		s.initMemory()
	}

	if interval > 0 {
		go func(ctx context.Context) {
			storeTicker := time.NewTicker(interval)
			for {
				select {
				case <-storeTicker.C:
					s.storing()
				case <-ctx.Done():
					storeTicker.Stop()
					return
				}
			}
		}(ctx)
	}
}

func (s *scheduler) storing() {
	s.repo.SaveMetrics(s.memory.ReadAll())
}

func (s *scheduler) initMemory() {
	s.memory.InitializeMetrics(s.repo.ReadMetrics())
}
