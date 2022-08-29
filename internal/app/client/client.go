package client

// Конструкторы для слоев и graceful shutdown

import (
	"context"
	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/controller/client/v1"
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

	v1.NewRouter(ctx, cfg, client, agentUseCase)
}
