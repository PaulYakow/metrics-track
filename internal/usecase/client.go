package usecase

import (
	"github.com/PaulYakow/metrics-track/internal/usecase/services/gather"
)

// Реализация клиента

type Client struct {
	repo   IClientRepo
	gather IClientGather
}

func NewClientUC(r IClientRepo) *Client {
	return &Client{repo: r, gather: gather.New()}
}

func (c *Client) Poll() {
	c.repo.Store(c.gather.Update())
}

func (c *Client) UpdateRoutes() []string {
	return c.repo.ReadCurrentMetrics()
}

func (c *Client) UpdateValues() [][]byte {
	return c.repo.ReadCurrentValues()
}
