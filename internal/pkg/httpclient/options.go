package httpclient

import "time"

type Option func(*Client)

func Address(address string) Option {
	return func(c *Client) {
		c.client.SetBaseURL(address)
	}
}

func Timeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.client.SetTimeout(timeout)
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.shutdownTimeout = timeout
	}
}
