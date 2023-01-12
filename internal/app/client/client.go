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

type Client struct {
	config    *config.Config
	logger    *logger.Logger
	repo      usecase.IClientMemory
	usecase   usecase.IClient
	hasher    usecase.IHasher
	collector *client.Collector
	sender    client.ISender
}

// New собирает клиента из слоёв (хранилище, логика, сервисы).
func New(cfg *config.Config) *Client {
	c := &Client{
		config: cfg,
		logger: logger.New(),
		repo:   repo.NewClientRepo(),
		hasher: hasher.New(cfg.Key),
	}

	c.usecase = usecase.NewClientUC(context.Background(), c.repo, c.hasher)
	c.collector = client.NewCollector(c.usecase, c.logger)

	if cfg.UseHTTPClient() {
		endpoint := fmt.Sprintf("http://%s/updates/", cfg.Address)
		c.sender = client.NewHTTPSender(httpclient.New(httpclient.RealIP(cfg.RealIP)), c.usecase, endpoint, c.logger, cfg)

	} else if cfg.UseGRPCClient() {
		c.sender = client.NewGRPCSender(c.usecase, c.logger)
	}

	return c
}

// Run запускает отдельными потоками "сборщика" метрик и отправку данных.
// В конце организован graceful shutdown.
func (c *Client) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go c.collector.Run(ctx, wg, c.config)
	go c.sender.Run(ctx, wg, c.config)

	// Ожидание сигнала завершения
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	s := <-interrupt
	cancel()
	wg.Wait()
	c.logger.Info("client - Run - signal: %s", s.String())
}
