package client

// Конструкторы для слоев и graceful shutdown

import (
	"context"
	"fmt"
	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/controller/client"
	"github.com/PaulYakow/metrics-track/internal/pkg/httpclient"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Run(cfg *config.ClientCfg) {
	ctx, cancel := context.WithCancel(context.Background())

	wg := new(sync.WaitGroup)

	l := logger.New()

	agentRepo := repo.NewClientRepo()

	agentUseCase := usecase.NewClientUC(agentRepo)

	collector := client.NewCollector(ctx, agentUseCase, l)
	wg.Add(1)
	go collector.Run(wg, cfg.PollInterval)

	c := httpclient.New(ctx)
	endpoint := fmt.Sprintf("http://%s/update/", cfg.Address)
	sender := client.NewSender(c, agentUseCase, endpoint, l)
	wg.Add(1)
	go sender.Run(wg, cfg.ReportInterval)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	s := <-interrupt
	cancel()
	wg.Wait()
	l.Info("client - Run - signal: %s", s.String())
}
