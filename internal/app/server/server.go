// Package server точка входа сервера сбора и хранения метрик.
package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PaulYakow/metrics-track/cmd/server/config"
	"github.com/PaulYakow/metrics-track/internal/controller/server"
	serverCtrl "github.com/PaulYakow/metrics-track/internal/controller/server/v1"
	v2 "github.com/PaulYakow/metrics-track/internal/controller/server/v2"
	"github.com/PaulYakow/metrics-track/internal/pkg/httpserver"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	postgre "github.com/PaulYakow/metrics-track/internal/pkg/postgre/v2"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/hasher"
)

type Server struct {
	config  *config.Config
	logger  *logger.Logger
	repo    usecase.IServerRepo
	hasher  usecase.IHasher
	usecase usecase.IServer
	httpSrv *httpserver.Server
	grpcSrv *v2.MetricsServer
}

// New собирает сервер из слоёв (хранилище, логика, сервисы).
func New(cfg *config.Config) *Server {
	s := &Server{
		config: cfg,
		logger: logger.New(),
		hasher: hasher.New(cfg.Key),
	}

	if cfg.Dsn != "" {
		s.repo = s.createDBRepo()
	} else if cfg.StoreFile != "" && s.repo == nil {
		s.repo = s.createMemoryRepo()
	}

	s.usecase = usecase.NewServerUC(s.repo, s.hasher)

	// HTTP server
	if cfg.UseHTTPServer() {
		handler := serverCtrl.NewRouter(s.usecase, s.logger, cfg)
		s.httpSrv = httpserver.New(handler, httpserver.Address(cfg.Address))
		s.logger.Info("server - run with params: a=%s | i=%v | f=%s | r=%v | k=%v | d=%s | crypto=%s",
			cfg.Address, cfg.StoreInterval, cfg.StoreFile, cfg.Restore, cfg.Key, cfg.Dsn, cfg.PathToCryptoKey)
	}

	// gRPC server
	if cfg.UseGRPCServer() {
		s.grpcSrv = v2.New(s.usecase, s.logger, cfg)
	}

	return s
}

// Run запускает сервер.
// В конце организован graceful shutdown.
func (s *Server) Run() {
	defer s.logger.Exit()

	if _, ok := s.repo.(*repo.ServerMemoryRepo); ok {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		s.startScheduler(ctx)
	}

	if s.httpSrv != nil {
		s.httpSrv.Run()
	}

	if s.grpcSrv != nil {
		s.grpcSrv.Run()
	}

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-interrupt
	s.logger.Info("server - Run - signal: %v", sig.String())

	// Shutdown
	if err := s.repo.CloseConnection(); err != nil {
		s.logger.Error(fmt.Errorf("server - Run - close connection to repo: %w", err))
	}

	if err := s.httpSrv.Shutdown(); err != nil {
		s.logger.Error(fmt.Errorf("server - Run - shutdown httpserver: %w", err))
	}
}

func (s *Server) createMemoryRepo() *repo.ServerMemoryRepo {
	return repo.NewServerMemory()
}

func (s *Server) createDBRepo() *repo.ServerSqlxImpl {
	pg, err := postgre.New(s.config.Dsn)
	if err != nil {
		s.logger.Fatal(fmt.Errorf("server - Run - postgre.New: %w", err))
	}
	//defer pg.Close()
	s.logger.Info("server - Run - PSQL connection ok")

	serverRepo, err := repo.NewSqlxImpl(pg)
	if err != nil {
		s.logger.Fatal(fmt.Errorf("server - Run - repo.New: %w", err))
	}
	s.logger.Info("server - Run - PSQL in use")

	return serverRepo
}

func (s *Server) startScheduler(ctx context.Context) {
	scheduler, err := server.NewScheduler(s.repo, s.config.StoreFile, s.logger)
	if err != nil {
		s.logger.Error(fmt.Errorf("server - run scheduler: %w", err))
	}
	scheduler.Run(ctx, s.config.Restore, s.config.StoreInterval)
	s.logger.Info("server - Run - file storage in use")
}
