// Package v1 содержит маршруты конечных точек и шаблон страницы с метриками.
package v1

import (
	"fmt"
	"html/template"
	"net/http"
	"net/netip"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/PaulYakow/metrics-track/cmd/server/config"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/utils/pki"
)

const (
	updateRoute      = "/update"
	batchUpdateRoute = "/updates"
	valueRoute       = "/value"
	pingRoute        = "/ping"

	templateName = "./web/templates/metrics_list.gohtml"
)

var tmpl *template.Template

type serverRoutes struct {
	uc            usecase.IServer
	logger        logger.ILogger
	decoder       *pki.Decryptor
	trustedSubnet netip.Prefix
}

// NewRouter формирует основной роутер для обработки запросов (на основе chi).
func NewRouter(uc usecase.IServer, l logger.ILogger, cfg *config.Config) chi.Router {
	s := &serverRoutes{
		uc:     uc,
		logger: l,
	}

	if cfg.PathToCryptoKey != "" {
		var err error
		s.decoder, err = pki.NewDecryptor(cfg.PathToCryptoKey)
		if err != nil {
			s.logger.Fatal(err)
		}
	}

	if cfg.TrustedSubnet != "" {
		var err error
		s.trustedSubnet, err = netip.ParsePrefix(cfg.TrustedSubnet)
		if err != nil {
			s.logger.Fatal(err)
		}
	}

	tmpl = template.Must(template.ParseFiles(templateName))

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(compressGzip)

	if cfg.TrustedSubnet != "" {
		mux.Use(s.checkRealIP)
	}

	mux.Get("/", s.listOfMetrics)
	mux.Get(pingRoute, s.pingDB)
	mux.Mount(updateRoute, s.createUpdateRoutes())
	mux.Mount(batchUpdateRoute, s.createBatchUpdateRoutes())
	mux.Mount(valueRoute, s.createValueRoutes())

	return mux
}

func (s *serverRoutes) listOfMetrics(rw http.ResponseWriter, r *http.Request) {
	data, err := s.uc.GetAll(r.Context())
	if err != nil {
		s.logger.Error(fmt.Errorf("router - GetAll metrics failed: %w", err))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(rw, data)
	if err != nil {
		s.logger.Error(fmt.Errorf("router - apply template failed: %w", err))
	}
}

func (s *serverRoutes) pingDB(rw http.ResponseWriter, r *http.Request) {
	err := s.uc.CheckRepo()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
