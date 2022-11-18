package client

// Конструкторы для слоев и graceful shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/controller/client"
	"github.com/PaulYakow/metrics-track/internal/pkg/httpclient"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/hasher"
)

func Run(cfg *config.ClientCfg) {
	ctx, cancel := context.WithCancel(context.Background())

	l := logger.New()

	agentRepo := repo.NewClientRepo()
	agentHasher := hasher.New(cfg.Key)

	agentUseCase := usecase.NewClientUC(ctx, agentRepo, agentHasher)

	collector := client.NewCollector(agentUseCase, l)

	c := httpclient.New()
	endpoint := fmt.Sprintf("http://%s/updates/", cfg.Address)
	sender := client.NewSender(c, agentUseCase, endpoint, l)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go collector.Run(ctx, wg, cfg.PollInterval)
	go sender.Run(ctx, wg, cfg.ReportInterval)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	s := <-interrupt
	cancel()
	wg.Wait()
	l.Info("client - Run - signal: %s", s.String())
}
