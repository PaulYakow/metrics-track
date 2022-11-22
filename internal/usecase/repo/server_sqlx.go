package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/postgre/v2"
)

// ServerSqlxImpl реализация репозитория сервера (usecase.IServerRepo). Хранение в БД Postgres (драйвер - sqlx).
type ServerSqlxImpl struct {
	*v2.Postgre
}

var stmtUpsertMetric *sqlx.NamedStmt

// NewSqlxImpl создаёт объект ServerSqlxImpl.
func NewSqlxImpl(pg *v2.Postgre) (*ServerSqlxImpl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pg.ExecContext(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("repo - NewSqlxImpl - create table failed: %w", err)
	}

	stmtUpsertMetric, err = pg.PrepareNamed(upsertMetric)
	if err != nil {
		return nil, fmt.Errorf("new db - stmt prepare: %w", err)
	}

	return &ServerSqlxImpl{pg}, nil
}

func (repo *ServerSqlxImpl) Store(metric *entity.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := stmtUpsertMetric.ExecContext(ctx, *metric)
	if err != nil {
		return fmt.Errorf("repo - Store - update in DB: %w", err)
	}
	return nil
}

func (repo *ServerSqlxImpl) StoreBatch(metrics []entity.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := repo.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("repo - start transaction: %w", err)
	}
	defer tx.Rollback()

	txStmt := tx.NamedStmtContext(ctx, stmtUpsertMetric)

	for _, metric := range metrics {
		if _, err = txStmt.ExecContext(ctx, metric); err != nil {
			return fmt.Errorf("repo - ExecContext: %w", err)
		}
	}

	return tx.Commit()
}

func (repo *ServerSqlxImpl) Read(ctx context.Context, metric entity.Metric) (*entity.Metric, error) {
	ctxInner, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var result entity.Metric
	if err := repo.GetContext(ctxInner, &result, selectMetricByIDAndType, metric.ID, metric.MType); err != nil {
		return nil, fmt.Errorf("repo - Read: %w", err)
	}

	return &result, nil
}

func (repo *ServerSqlxImpl) ReadAll(ctx context.Context) ([]entity.Metric, error) {
	ctxInner, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var result []entity.Metric
	if err := repo.SelectContext(ctxInner, &result, selectAllMetrics); err != nil {
		return nil, fmt.Errorf("repo - Read [%v]: %w", result, err)
	}

	return result, nil
}

func (repo *ServerSqlxImpl) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := repo.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
