package httpclient

import (
	"context"
	"github.com/imroc/req/v3"
	"time"
)

const (
	_defaultAddr            = "http://localhost:8080/"
	_defaultTimeout         = 1 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

type Client struct {
	client          *req.Client
	ctx             context.Context
	shutdownTimeout time.Duration
}

func New(ctx context.Context, opts ...Option) *Client {
	httpclient := req.C().
		SetTimeout(_defaultTimeout)

	c := &Client{
		client:          httpclient,
		ctx:             ctx,
		shutdownTimeout: _defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Client) PostByURL(route string) error {
	_, err := c.client.R().
		SetContext(c.ctx).
		SetHeader("Content-Type", "plain/text").
		Post(route)

	return err
}

func (c *Client) PostByJSON(route string, data []byte) error {
	_, err := c.client.R().
		SetContext(c.ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(route)

	return err
}
