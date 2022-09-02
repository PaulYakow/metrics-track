package usecase

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/gather"
)

// Реализация клиента

type client struct {
	repo   IClientRepo
	gather IClientGather
	hasher IHasher
}

func NewClientUC(r IClientRepo, h IHasher) *client {
	return &client{
		repo:   r,
		gather: gather.New(),
		hasher: h,
	}
}

func (c *client) Poll() {
	c.repo.Store(c.gather.Update())
}

func (c *client) GetAll() []entity.Metric {
	return c.hasher.ProcessBatch(c.repo.ReadAll())
}
