package usecase

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/gather"
)

// Реализация клиента

type client struct {
	repo     IClientMemory
	hasher   IHasher
	gatherRT IClientGather
	gatherPS IClientGather
}

func NewClientUC(ctx context.Context, r IClientMemory, h IHasher) *client {
	return &client{
		repo:     r,
		hasher:   h,
		gatherRT: gather.NewGatherRuntime(),
		gatherPS: gather.NewGatherPsutil(ctx),
	}
}

func (c *client) Poll() {
	c.repo.Store(c.gatherRT.Update())
	c.repo.Store(c.gatherPS.Update())
}

func (c *client) GetAll() []entity.Metric {
	return c.hasher.ProcessBatch(c.repo.ReadAll())
}
