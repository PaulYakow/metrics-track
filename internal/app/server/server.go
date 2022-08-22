package server

import (
	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/controller/httpserver"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"log"
	"net/http"
)

func Run(cfg *config.ServerCfg) {
	serverRepo := repo.NewServerRepo()

	serverUseCase := usecase.NewServerUC(serverRepo)

	log.Fatal(http.ListenAndServe(cfg.Address, httpserver.NewRouter(serverUseCase)))
}
