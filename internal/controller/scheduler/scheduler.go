package scheduler

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"time"
)

type scheduler struct {
	uc usecase.ISchedule
}

func NewScheduler(ctx context.Context, uc usecase.ISchedule, restore bool, interval time.Duration) {
	sched := &scheduler{
		uc: uc,
	}

	if restore {
		sched.uc.InitMetrics()
	}

	if interval > 0 {
		go func(ctx context.Context) {
			storeTicker := time.NewTicker(interval)
			for {
				select {
				case <-storeTicker.C:
					sched.uc.RunStoring()
				case <-ctx.Done():
					storeTicker.Stop()
					return
				}
			}
		}(ctx)
	}
}
