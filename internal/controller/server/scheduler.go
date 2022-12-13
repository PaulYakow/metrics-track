// Package server сервис периодического сохранения метрик в файл (и восстановления при запуске).
package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/consumer"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/producer"
)

// Scheduler управляет периодическим сохранением метрик в файл.
type Scheduler struct {
	repo     usecase.IServerRepo
	producer *producer.Producer
	consumer *consumer.Consumer
	logger   logger.ILogger
}

// NewScheduler создаёт объект Scheduler.
func NewScheduler(repo usecase.IServerRepo, filename string, l logger.ILogger) (*Scheduler, error) {
	if filename != "" {
		p, err := producer.NewProducer(filename)
		if err != nil {
			return nil, err
		}

		c, err := consumer.NewConsumer(filename)
		if err != nil {
			return nil, err
		}

		return &Scheduler{
			repo:     repo,
			producer: p,
			consumer: c,
			logger:   l,
		}, nil
	}

	return nil, nil
}

// Run - загружает начальные значения метрик при старте и запускает периодическое сохранение в файл.
func (s *Scheduler) Run(ctx context.Context, restore bool, interval time.Duration) {
	if restore {
		s.initMemory()
	}

	if interval > 0 {
		go func(ctx context.Context) {
			storeTicker := time.NewTicker(interval)
			for {
				select {
				case <-storeTicker.C:
					s.storing(ctx)
				case <-ctx.Done():
					storeTicker.Stop()
					return
				}
			}
		}(ctx)
	}
}

func (s *Scheduler) storing(ctx context.Context) {
	metrics, err := s.repo.ReadAll(ctx)
	if err != nil {
		s.logger.Error(fmt.Errorf("scheduler - read all metrics: %w", err))
	}

	err = s.producer.Write(&metrics)
	if err != nil {
		s.logger.Error(fmt.Errorf("scheduler - save to file: %w", err))
	}
}

func (s *Scheduler) initMemory() {
	metrics, err := s.consumer.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			s.logger.Info("scheduler - file is empty (no data for init)")
			return
		}

		s.logger.Error(fmt.Errorf("scheduler - read from file: %w", err))
		return
	}

	for _, metric := range metrics {
		if err = s.repo.Store(metric); err != nil {
			s.logger.Error(fmt.Errorf("scheduler - init metric: %w", err))
		}
	}
}
