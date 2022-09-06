package repo

import (
	"context"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/postgre"
	"log"
	"time"
)

const (
	_defaultEntityCap = 30
	_createTable      = `
CREATE TABLE IF NOT EXISTS metrics(
    "id" VARCHAR(255) UNIQUE NOT NULL,
    "type" VARCHAR(50) NOT NULL,
    "value" DOUBLE PRECISION,
    "delta" INTEGER,
    "hash" VARCHAR(255)
    );
`
	_upsertMetric = `
INSERT INTO metrics (id, type, value, delta, hash)
VALUES($1,$2,$3,$4,$5) 
ON CONFLICT (id) DO UPDATE
SET value = EXCLUDED.value, delta = EXCLUDED.delta, hash = EXCLUDED.hash;
`
	_readMetrics = `SELECT * FROM metrics`
	_readMetric  = `SELECT $1,$2,$3,$4,$5 FROM metrics`
)

type serverPSQLRepo struct {
	*postgre.Postgre
}

func NewServerPostgre(pg *postgre.Postgre) (*serverPSQLRepo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pg.Pool.Exec(ctx, _createTable)
	if err != nil {
		return nil, fmt.Errorf("repo - NewServerPostgre - create table failed: %w", err)
	}

	return &serverPSQLRepo{pg}, nil
}

func (repo *serverPSQLRepo) Store(metric *entity.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := repo.Pool.Exec(ctx, _upsertMetric, metric.ID, metric.MType, metric.Value, metric.Delta, metric.Hash)
	if err != nil {
		return fmt.Errorf("repo - save metric - update/insert: %w", err)
	}

	return nil
}

func (repo *serverPSQLRepo) Read(metric entity.Metric) (*entity.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := repo.Pool.QueryRow(ctx, _readMetric, metric.ID, metric.MType, metric.Value, metric.Delta, metric.Hash)
	m := &entity.Metric{}
	err := row.Scan(&m.ID, &m.MType, &m.Value, &m.Delta, &m.Hash)
	if err != nil {
		return nil, fmt.Errorf("repo - Read - row.Scan: %w", err)
	}

	return m, nil
}

func (repo *serverPSQLRepo) ReadAll() []entity.Metric {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := repo.Pool.Query(ctx, _readMetrics)
	if err != nil {
		log.Printf("repo - ReadMetrics - Pool.Query: %v", err)
		return nil
	}
	defer rows.Close()

	metrics := make([]entity.Metric, 0, _defaultEntityCap)
	for rows.Next() {
		m := entity.Metric{}
		err = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta, &m.Hash)
		if err != nil {
			log.Printf("repo - ReadMetrics - rows.Scan: %v", err)
			return nil
		}
		metrics = append(metrics, m)
	}

	return metrics
}

func (repo *serverPSQLRepo) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := repo.Pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}
