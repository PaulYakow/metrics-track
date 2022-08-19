package client

// Конструкторы для слоев и graceful shutdown

import (
	"context"
	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/controller/httpclient"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"github.com/imroc/req/v3"
	"time"
)

func Run(ctx context.Context, cfg *config.ClientCfg) {
	agentRepo := repo.NewClientRepo()

	agentUseCase := usecase.NewClientUC(agentRepo)

	client := req.C().
		SetTimeout(time.Second)

	httpclient.NewRouter(ctx, cfg, client, agentUseCase)
}
