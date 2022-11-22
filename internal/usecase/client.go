package usecase

import (
	"context"

	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/gather"
)

// Client реализация контроллера клиента (IClient).
type Client struct {
	repo     IClientMemory
	hasher   IHasher
	gatherRT IClientGather
	gatherPS IClientGather
}

// NewClientUC создаёт объект Client.
func NewClientUC(ctx context.Context, r IClientMemory, h IHasher) *Client {
	return &Client{
		repo:     r,
		hasher:   h,
		gatherRT: gather.NewGatherRuntime(),
		gatherPS: gather.NewGatherPsutil(ctx),
	}
}

func (c *Client) Poll() {
	c.repo.Store(c.gatherRT.Update())
	c.repo.Store(c.gatherPS.Update())
}

func (c *Client) GetAll() []entity.Metric {
	return c.hasher.ProcessBatch(c.repo.ReadAll())
}
