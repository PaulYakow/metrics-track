package server

import (
	"github.com/PaulYakow/metrics-track/internal/controller/httpserver"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"log"
	"net/http"
)

const (
	endpoint = ":8080"
)

func Run() {
	serverRepo := repo.NewServerRepo()

	serverUseCase := usecase.NewServerUC(serverRepo)

	log.Fatal(http.ListenAndServe(endpoint, httpserver.NewRouter(serverUseCase)))
}
