package v1

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PaulYakow/metrics-track/cmd/server/config"
	"github.com/PaulYakow/metrics-track/internal/pkg/httpserver"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/hasher"
)

func ExampleNewRouter() {
	var err error

	// Создаём логгер
	l := logger.New()
	defer l.Exit()

	// Объявляем нужную реализацию хранилища
	var someRepo usecase.IServerRepo = repo.NewServerMemory()

	// Также нужен сервис проверки хэша
	serverHasher := hasher.New("some key from config")

	// Слой логики, который связывает хранилище и сервисы
	serverUseCase := usecase.NewServerUC(someRepo, serverHasher)

	// И наконец HTTP-сервер, который обрабатывает запросы
	handler := NewRouter(serverUseCase, l, &config.Config{})
	srv := httpserver.New(handler, httpserver.Address(":8080"))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем один из сигналов прерывания (от системы либо ошибка сервера
	s := <-interrupt
	l.Info("server - Run - signal: %v", s.String())

	// Завершение сервера
	err = srv.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("server - Run - Shutdown: %w", err))
	}
}
