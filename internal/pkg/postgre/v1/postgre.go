package v1

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

const (
	_defaultMaxPoolSize     = 4
	_defaultMaxConnIdleTime = time.Second * 30
	_defaultMaxConnLifetime = time.Minute * 2

	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

//goland:noinspection SpellCheckingInspection
type Postgre struct {
	maxPollSize     int
	maxConnIdleTime time.Duration
	maxConnLifeTime time.Duration
	connAttempts    int
	connTimeout     time.Duration

	Pool *pgxpool.Pool
}

func New(dsn string, opts ...Option) (*Postgre, error) {
	pg := &Postgre{
		maxPollSize:     _defaultMaxPoolSize,
		maxConnIdleTime: _defaultMaxConnIdleTime,
		maxConnLifeTime: _defaultMaxConnLifetime,
		connAttempts:    _defaultConnAttempts,
		connTimeout:     _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("postgre - NewPostgre - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPollSize)
	poolConfig.MaxConnIdleTime = pg.maxConnIdleTime
	poolConfig.MaxConnLifetime = pg.maxConnLifeTime

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgre - NewPostgre - connAttempts == 0: %w", err)
	}

	return pg, nil
}

func (p *Postgre) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
