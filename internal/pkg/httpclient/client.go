package httpclient

import (
	"context"
	"time"

	"github.com/imroc/req/v3"
)

const (
	defaultAddr            = "http://localhost:8080/"
	defaultTimeout         = 1 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

type Client struct {
	client          *req.Client
	shutdownTimeout time.Duration
}

func New(opts ...Option) *Client {
	httpclient := req.C().
		SetTimeout(defaultTimeout)

	c := &Client{
		client:          httpclient,
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) PostByURL(route string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "plain/text").
		Post(route)

	return err
}

func (c *Client) PostByJSON(route string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(route)

	return err
}

func (c *Client) PostByJSONBatch(route string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := c.client.R().
		SetContext(ctx).
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(route)

	return err
}
