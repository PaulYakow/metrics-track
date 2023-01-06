// Package client точка входа клиента сбора и отправки метрик.
package client

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/PaulYakow/metrics-track/cmd/agent/config"
	"github.com/PaulYakow/metrics-track/internal/controller/client"
	"github.com/PaulYakow/metrics-track/internal/pkg/httpclient"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/hasher"
)

// Run собирает клиента из слоёв (хранилище, логика, сервисы).
// Запускает отдельными потоками "сборщика" метрик и отправку данных.
// В конце организован graceful shutdown.
func Run(cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())

	l := logger.New()

	agentRepo := repo.NewClientRepo()
	agentHasher := hasher.New(cfg.Key)

	agentUseCase := usecase.NewClientUC(ctx, agentRepo, agentHasher)

	collector := client.NewCollector(agentUseCase, l)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go collector.Run(ctx, wg, cfg)

	if cfg.Address != "" {
		go runHTTPSender(ctx, wg, agentUseCase, l, cfg)
	} else if cfg.Address == "" && cfg.GRPCTarget != "" {
		go runGRPCSender(ctx, wg, agentUseCase, l, cfg)
	}

	// Ожидание сигнала завершения
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	s := <-interrupt
	cancel()
	wg.Wait()
	l.Info("client - Run - signal: %s", s.String())
}

func runHTTPSender(ctx context.Context, wg *sync.WaitGroup, uc usecase.IClient, l logger.ILogger, cfg *config.Config) {
	c := httpclient.New(httpclient.RealIP(cfg.RealIP))
	endpoint := fmt.Sprintf("http://%s/updates/", cfg.Address)
	sender := client.NewHTTPSender(c, uc, endpoint, l, cfg)
	sender.Run(ctx, wg, cfg)
}

func runGRPCSender(ctx context.Context, wg *sync.WaitGroup, uc usecase.IClient, l logger.ILogger, cfg *config.Config) {
	sender := client.NewGRPCSender(uc, l)
	sender.Run(ctx, wg, cfg)
}
