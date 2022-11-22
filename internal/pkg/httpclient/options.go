package httpclient

import "time"

// Option применяет заданную настройку к клиенту (Client).
type Option func(*Client)

// Address задаёт клиенту базовый URL по умолчанию.
func Address(address string) Option {
	return func(c *Client) {
		c.client.SetBaseURL(address)
	}
}

// Timeout задаёт время ожидания для запросов.
func Timeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.client.SetTimeout(timeout)
	}
}

// ShutdownTimeout задаёт таймаут выключения
func ShutdownTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.shutdownTimeout = timeout
	}
}
