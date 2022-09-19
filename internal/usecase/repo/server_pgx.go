package repo

import (
	"context"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/postgre/v1"
	"time"
)

const (
	defaultEntityCap = 30
)

type serverPgxImpl struct {
	*v1.Postgre
}

func NewPgxImpl(pg *v1.Postgre) (*serverPgxImpl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pg.Pool.Exec(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("repo - NewPgxImpl - create table failed: %w", err)
	}

	return &serverPgxImpl{pg}, nil
}

func (repo *serverPgxImpl) Store(metric *entity.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := repo.Pool.Exec(ctx, upsertMetric, metric.ID, metric.MType, metric.Delta, metric.Value, metric.Hash)
	if err != nil {
		return fmt.Errorf("repo - Store - update in DB: %w", err)
	}

	return nil
}

func (repo *serverPgxImpl) Read(metric entity.Metric) (*entity.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	m := &entity.Metric{}
	//var hash sql.NullString
	err := repo.Pool.QueryRow(ctx, selectMetricByIDAndType, metric.ID, metric.MType).
		Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash)
	if err != nil {
		return nil, fmt.Errorf("repo - Read - row.Scan: %w", err)
	}

	//m.Hash = hash.String
	return m, nil
}

func (repo *serverPgxImpl) ReadAll() ([]entity.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := repo.Pool.Query(ctx, selectAllMetrics)
	if err != nil {
		return nil, fmt.Errorf("repo - ReadAll - Pool.Query: %w", err)
	}
	defer rows.Close()

	metrics := make([]entity.Metric, 0, defaultEntityCap)
	for rows.Next() {
		m := entity.Metric{}
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash)
		if err != nil {
			return nil, fmt.Errorf("repo - ReadAll - rows.Scan: %w", err)
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (repo *serverPgxImpl) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := repo.Pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}
