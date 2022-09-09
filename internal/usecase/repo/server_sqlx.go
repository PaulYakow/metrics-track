package repo

import (
	"context"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/postgre/v2"
	"github.com/jmoiron/sqlx"
	"time"
)

type serverSqlxImpl struct {
	*v2.Postgre
}

var stmtUpsertMetric *sqlx.NamedStmt

func NewSqlxImpl(pg *v2.Postgre) (*serverSqlxImpl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pg.ExecContext(ctx, _schema)
	if err != nil {
		return nil, fmt.Errorf("repo - NewSqlxImpl - create table failed: %w", err)
	}

	stmtUpsertMetric, err = pg.PrepareNamed(_upsertMetric)
	if err != nil {
		return nil, fmt.Errorf("new db - stmt prepare: %w", err)
	}

	return &serverSqlxImpl{pg}, nil
}

func (repo *serverSqlxImpl) Store(metric *entity.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := stmtUpsertMetric.ExecContext(ctx, *metric)
	if err != nil {
		return fmt.Errorf("repo - Store - update in DB: %w", err)
	}
	return nil
}

func (repo *serverSqlxImpl) StoreBatch(metrics []entity.Metric) error {
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

func (repo *serverSqlxImpl) Read(metric entity.Metric) (*entity.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var result entity.Metric
	if err := repo.GetContext(ctx, &result, _selectMetricByIDAndType, metric.ID, metric.MType); err != nil {
		return nil, fmt.Errorf("repo - Read: %w", err)
	}

	return &result, nil
}

func (repo *serverSqlxImpl) ReadAll() ([]entity.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var result []entity.Metric
	if err := repo.SelectContext(ctx, &result, _selectAllMetrics); err != nil {
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
