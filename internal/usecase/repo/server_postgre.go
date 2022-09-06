package repo

import (
	"context"
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
	_upsertMetrics = `
INSERT INTO metrics (id, type, value, delta, hash)
VALUES($1,$2,$3,$4,$5) 
ON CONFLICT (id) DO UPDATE
SET value = EXCLUDED.value, delta = EXCLUDED.delta, hash = EXCLUDED.hash;
`
	_readMetrics = `SELECT * FROM metrics`
)

type serverPostgre struct {
	*postgre.Postgre
}

func NewServerPostgre(pg *postgre.Postgre) *serverPostgre {
	return &serverPostgre{pg}
}

func (s *serverPostgre) SaveMetrics(metrics []entity.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.Pool.Exec(ctx, _createTable)
	if err != nil {
		log.Printf("repo - save metrics - create table failed: %v", err)
		return
	}

	for _, metric := range metrics {
		_, err = s.Pool.Exec(ctx, _upsertMetrics, metric.ID, metric.MType, metric.Value, metric.Delta, metric.Hash)
		if err != nil {
			log.Printf("repo - save metric - update/insert: %v", err)
		}
	}
}

func (s *serverPostgre) ReadMetrics() []*entity.Metric {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.Pool.Query(ctx, _readMetrics)
	if err != nil {
		log.Printf("repo - ReadMetrics - Pool.Query: %v", err)
		return nil
	}
	defer rows.Close()

	metrics := make([]*entity.Metric, 0, _defaultEntityCap)
	for rows.Next() {
		m := &entity.Metric{}
		err = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta, &m.Hash)
		if err != nil {
			log.Printf("repo - ReadMetrics - rows.Scan: %v", err)
			return nil
		}
		metrics = append(metrics, m)
	}

	return metrics
}

func (s *serverPostgre) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.Pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}
