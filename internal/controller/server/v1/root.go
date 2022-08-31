package v1

import (
	"compress/flate"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"log"
	"net/http"
)

type serverRoutes struct {
	uc     usecase.IServer
	logger logger.ILogger
}

func NewRouter(uc usecase.IServer, l logger.ILogger) chi.Router {
	s := &serverRoutes{
		uc:     uc,
		logger: l,
	}

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Compress(flate.BestCompression))

	mux.Get("/", s.listOfMetrics) // return HTML with all metrics
	mux.Mount("/update", updates{}.Routes(s))
	mux.Mount("/value", values{}.Routes(s))

	return mux
}

func (s *serverRoutes) listOfMetrics(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./web/templates/metrics_list.gohtml"))
	data := s.uc.GetAllMetrics()
	rw.Header().Set("Content-Type", "text/html")
	err := tmpl.Execute(rw, data)
	if err != nil {
		log.Printf("apply template: %v", err)
	}
}
