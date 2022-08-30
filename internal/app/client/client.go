package client

// Конструкторы для слоев и graceful shutdown

import (
	"context"
	"fmt"
	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/controller/client"
	"github.com/PaulYakow/metrics-track/internal/pkg/httpclient"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(ctx context.Context, cfg *config.ClientCfg) {
	agentRepo := repo.NewClientRepo()

	agentUseCase := usecase.NewClientUC(agentRepo)

	collector := client.NewCollector(ctx, agentUseCase)
	go collector.Run(cfg.PollInterval)

	c := httpclient.New(ctx)
	endpoint := fmt.Sprintf("http://%s/update/", cfg.Address)
	sender := client.NewSender(c, agentUseCase, endpoint, cfg.ReportInterval)
	go sender.Run()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	s := <-interrupt
	log.Printf("client - Run - signal: %s", s.String())
}
