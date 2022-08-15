package client

// Конструкторы для слоев и graceful shutdown

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/controller/httpclient"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"github.com/imroc/req/v3"
	"time"
)

func Run(ctx context.Context) {
	metricsRep := repo.New()

	metricsUseCase := usecase.New(metricsRep)

	agent := httpclient.New(metricsUseCase)

	client := req.C().
		SetTimeout(time.Second)

	agent.Run(ctx, client)
}
