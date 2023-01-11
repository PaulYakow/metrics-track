// Package v2 содержит реализацию gRPC-сервера.
package v2

import (
	"context"
	"fmt"
	"net"
	"net/netip"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/PaulYakow/metrics-track/cmd/server/config"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	pb "github.com/PaulYakow/metrics-track/proto"
)

// MetricsServer управляет обработкой gRPC запросов.
type MetricsServer struct {
	pb.UnimplementedMetricsServer
	address       string
	uc            usecase.IServer
	logger        logger.ILogger
	trustedSubnet netip.Prefix
}

// New создаёт объект типа MetricsServer и запускает gRPC-сервер.
func New(uc usecase.IServer, l logger.ILogger, cfg *config.Config) *MetricsServer {
	s := &MetricsServer{
		address: cfg.GRPCAddress,
		uc:      uc,
		logger:  l,
	}

	if cfg.TrustedSubnet != "" {
		var err error
		s.trustedSubnet, err = netip.ParsePrefix(cfg.TrustedSubnet)
		if err != nil {
			s.logger.Fatal(err)
		}
	}

	return s
}

func (s *MetricsServer) Run() {
	go func() {
		listen, err := net.Listen("tcp", s.address)
		if err != nil {
			s.logger.Fatal(fmt.Errorf("gRPC - net.Listen: %w", err))
		}

		// создаём gRPC-сервер без зарегистрированной службы
		grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(s.checkIPInterceptor))
		// регистрируем сервис
		pb.RegisterMetricsServer(grpcSrv, s)

		s.logger.Info("gRPC run: %s", s.address)

		// получаем запрос gRPC
		if err = grpcSrv.Serve(listen); err != nil {
			s.logger.Fatal(fmt.Errorf("gRPC - Serve: %w", err))
		}
	}()
}

// UpdateBatch - сохранение пакета метрик.
func (s *MetricsServer) UpdateBatch(ctx context.Context, in *pb.UpdateBatchRequest) (*pb.UpdateBatchResponse, error) {
	var response pb.UpdateBatchResponse
	grpcMetrics := in.GetMetrics()
	entityMetrics := make([]entity.Metric, len(grpcMetrics))

	for i, metric := range grpcMetrics {
		entityMetrics[i] = entity.Metric{
			ID:    metric.Id,
			MType: metric.MType,
			Hash:  entity.NullString(metric.Hash),
		}

		switch metric.MType {
		case "gauge":
			entityMetrics[i].Value = &metric.Value
		case "counter":
			entityMetrics[i].Delta = &metric.Delta
		default:
		}
	}

	s.logger.Info("gRPC - update batch: %v", entityMetrics)

	if err := s.uc.SaveBatch(entityMetrics); err != nil {
		s.logger.Error(fmt.Errorf("gRPC - update batch save to storage: %w", err))
		response.Error = fmt.Sprintf("update batch save to storage: %v", err)
	}

	return &response, nil
}

// UpdateSingle - сохранение одиночной метрики.
func (s *MetricsServer) UpdateSingle(ctx context.Context, in *pb.UpdateSingleRequest) (*pb.UpdateSingleResponse, error) {
	var response pb.UpdateSingleResponse
	grpcMetric := in.GetMetric()
	entityMetric := &entity.Metric{
		ID:    grpcMetric.Id,
		MType: grpcMetric.MType,
		Hash:  entity.NullString(grpcMetric.Hash),
	}

	switch grpcMetric.MType {
	case "gauge":
		entityMetric.Value = &grpcMetric.Value
	case "counter":
		entityMetric.Delta = &grpcMetric.Delta
	default:
	}

	s.logger.Info("gRPC - update single: %v", entityMetric)

	if err := s.uc.Save(entityMetric); err != nil {
		s.logger.Error(fmt.Errorf("gRPC - update single save to storage: %w", err))
		response.Error = fmt.Sprintf("update single save to storage: %v", err)
	}

	return &response, nil
}

// GetSingle - чтение одиночной метрики.
func (s *MetricsServer) GetSingle(ctx context.Context, in *pb.GetSingleRequest) (*pb.GetSingleResponse, error) {
	var response pb.GetSingleResponse
	reqMetric := in.GetMetric()
	entityMetric := entity.Metric{
		ID:    reqMetric.Id,
		MType: reqMetric.MType,
		Hash:  entity.NullString(reqMetric.Hash),
	}

	switch reqMetric.MType {
	case "gauge":
		entityMetric.Value = &reqMetric.Value
	case "counter":
		entityMetric.Delta = &reqMetric.Delta
	default:
	}

	respMetric, err := s.uc.Get(ctx, entityMetric)
	if err != nil {
		s.logger.Error(fmt.Errorf("gRPC - get single from storage: %w", err))
		response.Error = fmt.Sprintf("get single from storage: %v", err)
	} else {
		response.Metric = &pb.Metric{
			Id:    respMetric.ID,
			MType: respMetric.MType,
			Hash:  string(respMetric.Hash),
		}
		switch respMetric.MType {
		case "gauge":
			response.Metric.Value = *respMetric.Value
		case "counter":
			response.Metric.Delta = *respMetric.Delta
		default:
		}
	}

	return &response, nil
}

// ListMetrics - чтение всех известных на данный момент метрик.
func (s *MetricsServer) ListMetrics(ctx context.Context, in *pb.ListMetricsRequest) (*pb.ListMetricsResponse, error) {
	var response pb.ListMetricsResponse
	var grpcMetrics []*pb.Metric

	entityMetrics, err := s.uc.GetAll(ctx)
	if err != nil {
		s.logger.Error(fmt.Errorf("gRPC - get all metrics from storage: %w", err))
		response.Error = fmt.Sprintf("get all metrics from storage: %v", err)
	} else {
		grpcMetrics = make([]*pb.Metric, len(entityMetrics))
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
		response.Metrics = grpcMetrics
	}

	return &response, nil
}

// CheckRepo - проверка соединения с базой данных.
func (s *MetricsServer) CheckRepo(ctx context.Context, in *pb.CheckRepoRequest) (*pb.CheckRepoResponse, error) {
	var response pb.CheckRepoResponse

	err := s.uc.CheckRepo()
	if err != nil {
		s.logger.Error(fmt.Errorf("gRPC - no connection to storage: %w", err))
		response.Error = fmt.Sprintf("no connection to storage: %v", err)
	}

	return &response, nil
}

func (s *MetricsServer) checkIPInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var realIP string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("real_ip")
		if len(values) > 0 {
			realIP = values[0]
		}
	}

	if realIP == "" {
		s.logger.Error(fmt.Errorf("missing trusted subnet in metadata"))
		return nil, status.Error(codes.FailedPrecondition, "missing trusted subnet")
	}

	ip, err := netip.ParseAddr(realIP)
	if err != nil {
		s.logger.Error(fmt.Errorf("check IP of agent: %w", err))
		return nil, status.Error(codes.Internal, "check IP error")
	}

	if !s.trustedSubnet.Contains(ip) {
		s.logger.Error(fmt.Errorf("no such IP in trusted: %s", ip))
		return nil, status.Error(codes.PermissionDenied, "no such IP in trusted")
	}

	return handler(ctx, req)
}
