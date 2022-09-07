package v1

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/consumer"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/producer"
	"log"
	"time"
)

type scheduler struct {
	repo     usecase.IServerRepo
	producer *producer.Producer
	consumer *consumer.Consumer
}

func NewScheduler(repo usecase.IServerRepo, filename string) (*scheduler, error) {
	if filename != "" {
		p, err := producer.NewProducer(filename)
		if err != nil {
			return nil, err
		}

		c, err := consumer.NewConsumer(filename)
		if err != nil {
			return nil, err
		}

		return &scheduler{
			repo:     repo,
			producer: p,
			consumer: c,
		}, nil
	}

	return nil, nil
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
	metrics, err := s.repo.ReadAll()
	if err != nil {
		log.Printf("scheduler - read all metrics: %v", err)
	}

	err = s.producer.Write(&metrics)
	if err != nil {
		log.Printf("scheduler - save in file: %v", err)
	}
}

func (s *scheduler) initMemory() {
	metrics, err := s.consumer.Read()
	if err != nil {
		log.Printf("scheduler - read from file: %v", err)
	}

	for _, metric := range metrics {
		if err = s.repo.Store(metric); err != nil {
			log.Printf("scheduler - init metric: %v", err)
		}
	}
}
