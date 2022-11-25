package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/postgre/v1"
)

const (
	defaultEntityCap = 30
)

// ServerPgxImpl реализация репозитория сервера (usecase.IServerRepo). Хранение в БД Postgres (драйвер - pgx).
type ServerPgxImpl struct {
	*v1.Postgre
}

// NewPgxImpl создаёт объект ServerPgxImpl.
func NewPgxImpl(pg *v1.Postgre) (*ServerPgxImpl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pg.Pool.Exec(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("repo - NewPgxImpl - create table failed: %w", err)
	}

	return &ServerPgxImpl{pg}, nil
}

func (repo *ServerPgxImpl) Store(metric *entity.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := repo.Pool.Exec(ctx, upsertMetric, metric.ID, metric.MType, metric.Delta, metric.Value, metric.Hash)
	if err != nil {
		return fmt.Errorf("repo - Store - update in DB: %w", err)
	}

	return nil
}

func (repo *ServerPgxImpl) Read(metric entity.Metric) (*entity.Metric, error) {
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

func (repo *ServerPgxImpl) ReadAll() ([]entity.Metric, error) {
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

func (repo *ServerPgxImpl) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := repo.Pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}
