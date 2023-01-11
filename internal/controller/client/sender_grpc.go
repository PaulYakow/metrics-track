// Package client содержит реализацию gRPC-клиента.
package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/PaulYakow/metrics-track/cmd/agent/config"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	pb "github.com/PaulYakow/metrics-track/proto"
)

const (
	defaultCtxTimeout = 500 * time.Millisecond
)

// GRPCSender управляет периодической отправкой метрик на заданный адрес.
type GRPCSender struct {
	client pb.MetricsClient
	uc     usecase.IClient
	logger logger.ILogger
}

// NewGRPCSender создаёт объект GRPCSender.
func NewGRPCSender(uc usecase.IClient, l logger.ILogger) *GRPCSender {
	s := &GRPCSender{
		uc:     uc,
		logger: l,
	}

	return s
}

// Run - запуск периодической отправки метрик.
func (s *GRPCSender) Run(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config) {
	ticker := time.NewTicker(cfg.ReportInterval)
	defer wg.Done()

	conn, err := grpc.Dial(cfg.GRPCTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.logger.Fatal(err)
	}
	defer conn.Close()

	s.client = pb.NewMetricsClient(conn)

	md := metadata.New(map[string]string{"real_ip": cfg.RealIP})
	ctx = metadata.NewOutgoingContext(ctx, md)

	s.logger.Info("gRPC sender - run with params: target=%s | report=%v | crypto=%s",
		cfg.GRPCTarget, cfg.ReportInterval, cfg.PathToCryptoKey)

	for {
		select {
		case <-ticker.C:
			s.sendMetrics(ctx)
		case <-ctx.Done():
			ticker.Stop()
			s.logger.Info("sender - context canceled")
			return
		}
	}
}

// sendMetrics - пакетная отправка метрик.
func (s *GRPCSender) sendMetrics(ctx context.Context) {
	entityMetrics := s.uc.GetAll()
	grpcMetrics := make([]*pb.Metric, len(entityMetrics))

	s.logger.Info("gRPC metrics: \n%v\n", entityMetrics)

	for i, metric := range entityMetrics {
		grpcMetrics[i] = &pb.Metric{
			Id:    metric.ID,
			MType: metric.MType,
			Hash:  string(metric.Hash),
		}

		switch metric.MType {
		case "gauge":
			grpcMetrics[i].Value = *metric.Value
		case "counter":
			grpcMetrics[i].Delta = *metric.Delta
		default:
		}
	}

	request := &pb.UpdateBatchRequest{Metrics: grpcMetrics}

	ctx, cancel := context.WithTimeout(ctx, defaultCtxTimeout)
	defer cancel()

	response, err := s.client.UpdateBatch(ctx, request)
	if err != nil {
		s.logger.Error(fmt.Errorf("gRPC sender - send batch of metrics: %w", err))
		return
	}
	s.logger.Info("gRPC sender - send batch of metrics: ", response.GetError())
}
