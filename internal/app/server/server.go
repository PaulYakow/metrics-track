// Package server точка входа сервера сбора и хранения метрик.
package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PaulYakow/metrics-track/cmd/server/config"
	"github.com/PaulYakow/metrics-track/internal/controller/server"
	serverCtrl "github.com/PaulYakow/metrics-track/internal/controller/server/v1"
	v2 "github.com/PaulYakow/metrics-track/internal/controller/server/v2"
	"github.com/PaulYakow/metrics-track/internal/pkg/httpserver"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	postgre "github.com/PaulYakow/metrics-track/internal/pkg/postgre/v2"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/hasher"
)

// Run собирает сервер из слоёв (хранилище, логика, сервисы).
// В конце организован graceful shutdown.
func Run(cfg *config.Config) {
	var err error
	l := logger.New()
	defer l.Exit()

	// In-memory storage
	var serverRepo usecase.IServerRepo = repo.NewServerMemory()

	serverHasher := hasher.New(cfg.Key)

	// File or db storage
	storage := false

	if cfg.Dsn != "" {
		pg, err1 := postgre.New(cfg.Dsn)
		if err1 != nil {
			l.Fatal(fmt.Errorf("server - Run - postgre.New: %w", err1))
		}
		defer pg.Close()
		l.Info("server - Run - PSQL connection ok")

		serverRepo, err1 = repo.NewSqlxImpl(pg)
		if err1 != nil {
			l.Fatal(fmt.Errorf("server - Run - repo.New: %w", err1))
		}
		storage = true
		l.Info("server - Run - PSQL in use")
	}

	if cfg.StoreFile != "" && !storage {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// memory <-> repo
		scheduler, err2 := server.NewScheduler(serverRepo, cfg.StoreFile, l)
		if err2 != nil {
			l.Error(fmt.Errorf("server - run scheduler: %w", err2))
		}
		scheduler.Run(ctx, cfg.Restore, cfg.StoreInterval)
		l.Info("server - Run - file storage in use")
	}

	serverUseCase := usecase.NewServerUC(serverRepo, serverHasher)

	// HTTP server
	handler := serverCtrl.NewRouter(serverUseCase, l, cfg)
	srv := httpserver.New(handler, httpserver.Address(cfg.Address))

	l.Info("server - run with params: a=%s | i=%v | f=%s | r=%v | k=%v | d=%s | crypto=%s",
		cfg.Address, cfg.StoreInterval, cfg.StoreFile, cfg.Restore, cfg.Key, cfg.Dsn, cfg.PathToCryptoKey)

	// gRPC server
	v2.New(serverUseCase, l, cfg)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-interrupt:
		l.Info("server - Run - signal: %v", s.String())
	case err = <-srv.Notify():
		l.Error(fmt.Errorf("server - Run - Notify: %w", err))
	}

	// Shutdown
	err = srv.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("server - Run - Shutdown: %w", err))
	}
}
