package v2

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

const (
	_defaultMaxOpenConn     = 4
	_defaultMaxIdleConn     = 4
	_defaultMaxConnIdleTime = time.Second * 30
	_defaultMaxConnLifetime = time.Minute * 2

	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

//goland:noinspection SpellCheckingInspection
type Postgre struct {
	maxOpenConn     int
	maxIdleConn     int
	maxConnIdleTime time.Duration
	maxConnLifetime time.Duration
	connAttempts    int
	connTimeout     time.Duration

	*sqlx.DB
}

func New(dsn string, opts ...Option) (*Postgre, error) {
	pg := &Postgre{
		maxOpenConn:     _defaultMaxOpenConn,
		maxIdleConn:     _defaultMaxIdleConn,
		maxConnIdleTime: _defaultMaxConnIdleTime,
		maxConnLifetime: _defaultMaxConnLifetime,
		connAttempts:    _defaultConnAttempts,
		connTimeout:     _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	var err error

	for pg.connAttempts > 0 {
		if pg.DB, err = sqlx.Connect("pgx", dsn); err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgre - NewPostgre - connAttempts == 0: %w", err)
	}

	pg.SetMaxOpenConns(pg.maxOpenConn)
	pg.SetMaxIdleConns(pg.maxIdleConn)
	pg.SetConnMaxIdleTime(pg.maxConnIdleTime)
	pg.SetConnMaxLifetime(pg.maxConnLifetime)

	return pg, nil
}
