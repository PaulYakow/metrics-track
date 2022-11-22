package httpserver

import (
	"net"
	"time"
)

// Option применяет заданную настройку к серверу (Server).
type Option func(*Server)

// Address задаёт серверу адрес для "прослушивания" (в формате "host:port")
func Address(address string) Option {
	return func(s *Server) {
		s.server.Addr = address
	}
}

// Port задаёт серверу порт для "прослушивания"
func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

// ReadTimeout задаёт максимальную продолжительность чтения запроса
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

// WriteTimeout задаёт максимальный таймаут при записи ответа
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

// ShutdownTimeout задаёт выдержку времени при завершении работы сервера
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
