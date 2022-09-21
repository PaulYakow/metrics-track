package v1

import (
	"compress/flate"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"net/http"
)

const (
	updateRoute      = "/update"
	batchUpdateRoute = "/updates"
	valueRoute       = "/value"
	pingRoute        = "/ping"

	templateName = "./web/templates/metrics_list.gohtml"
)

type serverRoutes struct {
	uc     usecase.IServer
	logger logger.ILogger
	tmpl   *template.Template
}

func NewRouter(uc usecase.IServer, l logger.ILogger) chi.Router {
	s := &serverRoutes{
		uc:     uc,
		logger: l,
		tmpl:   template.Must(template.ParseFiles(templateName)),
	}

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Compress(flate.BestCompression))

	mux.Get("/", s.listOfMetrics) // return HTML with all metrics
	mux.Mount(updateRoute, s.createUpdateRoutes())
	mux.Mount(batchUpdateRoute, s.createBatchUpdateRoutes())
	mux.Mount(valueRoute, s.createValueRoutes())
	mux.Mount(pingRoute, s.createPingRoutes())

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
	err = s.tmpl.Execute(rw, data)
	if err != nil {
		s.logger.Error(fmt.Errorf("router - apply template failed: %w", err))
	}
}
