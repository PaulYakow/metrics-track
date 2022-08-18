package usecase

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/gather"
)

// Адаптеры для клиента

type IClient interface {
	Poll()
	UpdateRoutes() []string
	UpdateValues() [][]byte
}

type IClientRepo interface {
	Store(map[string]*entity.Gauge, map[string]*entity.Counter)
	ReadCurrentMetrics() []string
	ReadCurrentValues() [][]byte
}

type IClientGather interface {
	Update() (map[string]*entity.Gauge, map[string]*entity.Counter)
}

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
