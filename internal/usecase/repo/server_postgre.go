package repo

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/postgre"
	"time"
)

type serverPostgre struct {
	*postgre.Postgre
}

func NewServerPostgre(pg *postgre.Postgre) *serverPostgre {
	return &serverPostgre{pg}
}

func (s serverPostgre) SaveMetrics(metrics []entity.Metric) {
	//TODO implement me
}

func (s serverPostgre) ReadMetrics() []*entity.Metric {
	//TODO implement me
	return nil
}

func (s serverPostgre) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.Pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}
