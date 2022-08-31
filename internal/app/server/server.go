package server

import (
	"context"
	"fmt"
	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/controller/scheduler"
	"github.com/PaulYakow/metrics-track/internal/controller/server/v1"
	"github.com/PaulYakow/metrics-track/internal/pkg/httpserver"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.ServerCfg) {
	l := logger.New()
	// In-memory repository
	serverMemory := repo.NewServerMemory()
	serverUseCase := usecase.NewServerUC(serverMemory)

	// File repository
	if cfg.StoreFile != "" {
		serverFile, err := repo.NewServerFile(cfg.StoreFile)
		if err != nil {
			l.Error(fmt.Errorf("server - create file storage: %v", err))
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		schedulerUseCase := usecase.NewScheduleUC(serverFile, serverMemory)
		scheduler.NewScheduler(ctx, schedulerUseCase, cfg.Restore, cfg.StoreInterval)
	}

	// HTTP server
	handler := v1.NewRouter(serverUseCase, l)
	server := httpserver.New(handler, httpserver.Address(cfg.Address))

	l.Info("server - run with params: a=%s | i=%v | f=%s | r=%v",
		cfg.Address, cfg.StoreInterval, cfg.StoreFile, cfg.Restore)

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
	err := server.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("server - Run - Shutdown: %w", err))
	}
}
