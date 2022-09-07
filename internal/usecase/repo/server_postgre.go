package repo

// todo: необходимо убрать лишние зависимости
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/postgre"
	"github.com/jackc/pgx/v4"
	"time"
)

const (
	_defaultEntityCap = 30
	_createTable      = `
CREATE TABLE IF NOT EXISTS metrics(
    "id" VARCHAR(255) UNIQUE NOT NULL,
    "type" VARCHAR(50) NOT NULL,
    "delta" INTEGER,
    "value" DOUBLE PRECISION,
    "hash" VARCHAR(255)
    );
`
	_upsertMetric = `
INSERT INTO metrics (id, type, delta, value, hash)
VALUES($1,$2,$3,$4,$5) 
ON CONFLICT (id) DO UPDATE
SET delta = EXCLUDED.delta, value = EXCLUDED.value, hash = EXCLUDED.hash;
`
	_readMetrics = `SELECT * FROM metrics;`
	_readMetric  = `
SELECT *
FROM metrics
WHERE id = $1 AND type = $2;
`
	_createRow = `
INSERT INTO metrics (id, type, delta, value, hash)
SELECT $1::VARCHAR,$2,$3,$4,$5
WHERE NOT EXISTS (
    SELECT 1 FROM metrics WHERE id = $1
)
RETURNING *;
`
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

	m, err := repo.Read(*metric)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err = repo.Pool.Exec(ctx, _createRow, metric.ID, metric.MType, metric.Delta, metric.Value, metric.Hash)
			if err != nil {
				return fmt.Errorf("repo - Store - try create row: %w", err)
			}
			return nil
		} else {
			return fmt.Errorf("repo - Store - not exists: %w", err)
		}
	}

	err = m.Update(metric)
	if err != nil {
		return fmt.Errorf("repo - Store - update metric: %w", err)
	}

	_, err = repo.Pool.Exec(ctx, _upsertMetric, m.ID, m.MType, m.Delta, m.Value, m.Hash)
	if err != nil {
		return fmt.Errorf("repo - Store - update in DB: %w", err)
	}

	return nil
}

func (repo *serverPSQLRepo) Read(metric entity.Metric) (*entity.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	m := &entity.Metric{}
	var hash sql.NullString
	err := repo.Pool.QueryRow(ctx, _readMetric, metric.ID, metric.MType).
		Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &hash)
	if err != nil {
		return nil, fmt.Errorf("repo - Read - row.Scan: %w", err)
	}

	m.Hash = hash.String
	return m, nil
}

func (repo *serverPSQLRepo) ReadAll() ([]entity.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := repo.Pool.Query(ctx, _readMetrics)
	if err != nil {
		return nil, fmt.Errorf("repo - ReadAll - Pool.Query: %w", err)
	}
	defer rows.Close()

	metrics := make([]entity.Metric, 0, _defaultEntityCap)
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

func (repo *serverPSQLRepo) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := repo.Pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}
