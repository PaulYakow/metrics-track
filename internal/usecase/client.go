package usecase

import (
	"context"

	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/gather"
)

// ClientUC реализация контроллера клиента (IClient).
type ClientUC struct {
	repo     IClientMemory
	hasher   IHasher
	gatherRT IClientGather
	gatherPS IClientGather
}

// NewClientUC создаёт объект ClientUC.
func NewClientUC(ctx context.Context, r IClientMemory, h IHasher) *ClientUC {
	return &ClientUC{
		repo:     r,
		hasher:   h,
		gatherRT: gather.NewGatherRuntime(),
		gatherPS: gather.NewGatherPsutil(ctx),
	}
}

func (c *ClientUC) Poll() {
	c.repo.Store(c.gatherRT.Update())
	c.repo.Store(c.gatherPS.Update())
}

func (c *ClientUC) GetAll() []entity.Metric {
	return c.hasher.ProcessBatch(c.repo.ReadAll())
}
