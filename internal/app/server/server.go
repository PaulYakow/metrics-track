package server

import (
	"context"
	"fmt"
	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/controller/server/v1"
	"github.com/PaulYakow/metrics-track/internal/pkg/httpserver"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/pkg/postgre"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/hasher"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.ServerCfg) {
	var err error
	l := logger.New()

	// In-memory storage
	var serverRepo usecase.IServerRepo = repo.NewServerMemory()

	serverHasher := hasher.New(cfg.Key)

	// File or db storage
	storage := false

	if cfg.Dsn != "" {
		pg, err := postgre.New(cfg.Dsn, postgre.MaxPoolSize(2))
		if err != nil {
			l.Fatal(fmt.Errorf("server - Run - postgre.New: %w", err))
		}
		defer pg.Close()

		serverRepo, err = repo.NewServerPostgre(pg)
		if err != nil {
			l.Fatal(fmt.Errorf("server - Run - repo.New: %w", err))
		}
		storage = true
	}

	if cfg.StoreFile != "" && !storage {
		// Server scheduler (memory <-> repo)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		scheduler, err := v1.NewScheduler(serverRepo, cfg.StoreFile)
		if err != nil {
			l.Error(fmt.Errorf("server - run scheduler: %w", err))
		}
		scheduler.Run(ctx, cfg.Restore, cfg.StoreInterval)
	}

	serverUseCase := usecase.NewServerUC(serverRepo, serverHasher)

	// HTTP server
	handler := v1.NewRouter(serverUseCase, l)
	server := httpserver.New(handler, httpserver.Address(cfg.Address))

	l.Info("server - run with params: a=%s | i=%v | f=%s | r=%v | d=%s",
		cfg.Address, cfg.StoreInterval, cfg.StoreFile, cfg.Restore, cfg.Dsn)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("server - Run - signal: %v", s.String())
	case err := <-server.Notify():
		l.Error(fmt.Errorf("server - Run - Notify: %w", err))
	}

	// Shutdown
	err = server.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("server - Run - Shutdown: %w", err))
	}
}
