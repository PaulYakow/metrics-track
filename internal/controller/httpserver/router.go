package httpserver

import (
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"
	"strings"
)

type serverRoutes struct {
	uc usecase.ServerMetric
}

func NewRouter(uc usecase.ServerMetric) chi.Router {
	s := &serverRoutes{uc: uc}

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Route("/", func(r chi.Router) {
		//Обработка GET-запроса к хосту
		r.Get("/", s.getListOfMetrics)

		r.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}", s.getMetricValue)
		})

		r.Route("/update", func(r chi.Router) {

			r.Route("/gauge", func(r chi.Router) {
				r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
					rw.WriteHeader(http.StatusNotFound)
				})
				r.Route("/{name}", func(r chi.Router) {
					r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusBadRequest)
					})
					r.Post("/{value}", s.postGauge)
				})
			})

			r.Route("/counter", func(r chi.Router) {
				r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
					rw.WriteHeader(http.StatusNotFound)
				})
				r.Route("/{name}", func(r chi.Router) {
					r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusBadRequest)
					})
					r.Post("/{value}", s.postCounter)
				})
			})

			r.Route("/", func(r chi.Router) {
				r.Post("/*", s.postDefault)
			})
		})
	})

	return mux
}

func (s *serverRoutes) postGauge(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	rawValue := chi.URLParam(r, "value")

	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	s.uc.SaveGauge(name, value)
	rw.WriteHeader(http.StatusOK)
}

func (s *serverRoutes) postCounter(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	rawValue := chi.URLParam(r, "value")

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	s.uc.SaveCounter(name, value)
	rw.WriteHeader(http.StatusOK)
}

func (s *serverRoutes) postDefault(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotImplemented)
}

func (s *serverRoutes) getListOfMetrics(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte(strings.Join(s.uc.GetAllMetrics(), "\n")))
}

func (s *serverRoutes) getMetricValue(rw http.ResponseWriter, r *http.Request) {
	typeName := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	value, err := s.uc.GetValueByType(typeName, name)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.Write([]byte(value))
}
