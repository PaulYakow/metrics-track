// Package httpclient содержит простейший http-клиент для отправки запросов посредством URL и JSON.
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

// Client http-клиент.
type Client struct {
	client          *req.Client
	shutdownTimeout time.Duration
}

// New - создаёт объект Client и применяет заданные настройки.
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

// PostByURL - отправка POST-запроса (Content-Type=plain/text) на заданный адрес.
func (c *Client) PostByURL(route string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "plain/text").
		Post(route)

	return err
}

// PostByJSON - отправка POST-запроса (Content-Type=application/json) на заданный адрес.
// data представляет собой одинарный JSON.
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

// PostByJSONBatch - отправка POST-запроса (Content-Type=application/json) на заданный адрес.
// data представляет собой пакет (массив) JSON.
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
