package server

import (
	"context"
	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/controller/scheduler"
	"github.com/PaulYakow/metrics-track/internal/controller/server/v1"
	"github.com/PaulYakow/metrics-track/internal/pkg/httpserver"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.ServerCfg) {
	// In-memory repository
	serverMemory := repo.NewServerMemory()
	serverUseCase := usecase.NewServerUC(serverMemory)

	// File repository
	if cfg.StoreFile != "" {
		serverFile, err := repo.NewServerFile(cfg.StoreFile)
		if err != nil {
			log.Printf("create file storage for server: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		schedulerUseCase := usecase.NewScheduleUC(serverFile, serverMemory)
		scheduler.NewScheduler(ctx, schedulerUseCase, cfg.Restore, cfg.StoreInterval)
	}

	// HTTP server
	handler := v1.NewRouter(serverUseCase)
	server := httpserver.New(handler, httpserver.Address(cfg.Address))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Printf("server - Run - signal: %s", s.String())
	case err := <-server.Notify():
		log.Printf("server - Run - Notify: %s", err)
	}

	// Shutdown
	err := server.Shutdown()
	if err != nil {
		log.Printf("server - Run - Shutdown: %s", err)
	}
}
