package server

import (
	"context"
	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/controller/httpserver"
	"github.com/PaulYakow/metrics-track/internal/controller/scheduler"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"log"
	"net/http"
)

func Run(cfg *config.ServerCfg) {
	serverMemory := repo.NewServerMemory()
	serverUseCase := usecase.NewServerUC(serverMemory)

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

	log.Fatal(http.ListenAndServe(cfg.Address, httpserver.NewRouter(serverUseCase)))
}
