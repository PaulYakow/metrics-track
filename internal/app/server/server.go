package server

import (
	"fmt"
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
	endpoint := fmt.Sprintf("%s:%s", cfg.Address[0], cfg.Address[1])
	log.Fatal(http.ListenAndServe(endpoint, httpserver.NewRouter(serverUseCase)))
}
