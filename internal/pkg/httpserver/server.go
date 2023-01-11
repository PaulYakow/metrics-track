// Package httpserver содержит простейший http-сервер для обработки запросов.
package httpserver

import (
	"context"
	"log"
	"net/http"
	"time"
)

const (
	defaultAddr            = ":8080"
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

// Server обёртка для http.Server
type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

// New - создаёт объект Server, применяет заданные настройки и запускает http-сервер на заданном порту.
func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		Addr:         defaultAddr,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	s := &Server{
		server:          httpServer,
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Shutdown - завершает работу сервера с выдержкой времени.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s *Server) Run() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err, s.Shutdown())
		}
	}()
}
