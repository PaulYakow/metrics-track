package v1

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
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
	uc     usecase.IServer
	logger logger.ILogger
}

func NewRouter(uc usecase.IServer, l logger.ILogger) chi.Router {
	s := &serverRoutes{
		uc:     uc,
		logger: l,
	}

	tmpl = template.Must(template.ParseFiles(templateName))

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(compressGzip)

	mux.Get("/", s.listOfMetrics) // return HTML with all metrics
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
