package repo

import (
	"context"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/postgre/v2"
	"time"
)

type serverSqlxImpl struct {
	*v2.Postgre
}

func NewSqlxImpl(pg *v2.Postgre) (*serverSqlxImpl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pg.ExecContext(ctx, _schema)
	if err != nil {
		return nil, fmt.Errorf("repo - NewSqlxImpl - create table failed: %w", err)
	}

	return &serverSqlxImpl{pg}, nil
}

func (repo *serverSqlxImpl) Store(metric *entity.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	m, err := repo.Read(*metric)
	if err != nil {
		_, err = repo.ExecContext(ctx, _createRow, metric.ID, metric.MType, metric.Delta, metric.Value, metric.Hash)
		if err != nil {
			return fmt.Errorf("repo - Store - try create row: %w", err)
		}
		return nil
	}

	err = m.Update(metric)
	if err != nil {
		return fmt.Errorf("repo - Store - update metric: %w", err)
	}

	_, err = repo.ExecContext(ctx, _upsertMetric, m.ID, m.MType, m.Delta, m.Value, m.Hash)
	if err != nil {
		return fmt.Errorf("repo - Store - update in DB: %w", err)
	}

	return nil
}

func (repo *serverSqlxImpl) Read(metric entity.Metric) (*entity.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var result entity.Metric
	if err := repo.GetContext(ctx, &result, _readMetric, metric.ID, metric.MType); err != nil {
		return nil, fmt.Errorf("repo - Read: %w", err)
	}

	return &result, nil
}

func (repo *serverSqlxImpl) ReadAll() ([]entity.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var result []entity.Metric
	if err := repo.SelectContext(ctx, &result, _readMetrics); err != nil {
		return nil, fmt.Errorf("repo - Read [%v]: %w", result, err)
	}

	return result, nil
}

func (repo *serverSqlxImpl) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := repo.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
