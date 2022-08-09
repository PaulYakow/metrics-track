package app

import (
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Server struct {
	sync.Mutex
	metrics map[string]model.Metric
}

func NewServer() *Server {
	return &Server{metrics: make(map[string]model.Metric)}
}

func (s *Server) Run() {
	router := s.newRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func (s *Server) postGaugeHandler(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	rawValue := chi.URLParam(r, "value")

	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	s.Lock()
	defer s.Unlock()
	if _, ok := s.metrics[name]; !ok {
		s.metrics[name] = &model.Gauge{}
	}

	s.metrics[name].SetValue(value)

	rw.WriteHeader(http.StatusOK)
}

func (s *Server) postCounterHandler(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	rawValue := chi.URLParam(r, "value")

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	s.Lock()
	defer s.Unlock()
	if _, ok := s.metrics[name]; !ok {
		s.metrics[name] = &model.Counter{}
	}

	oldValue := s.metrics[name].GetValue().(int64)
	s.metrics[name].SetValue(oldValue + int64(value))

	rw.WriteHeader(http.StatusOK)
}

func postDefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) getListOfMetrics(rw http.ResponseWriter, r *http.Request) {
	listOfMetrics := make([]string, 0)

	s.Lock()
	defer s.Unlock()
	for name, metric := range s.metrics {
		listOfMetrics = append(listOfMetrics,
			fmt.Sprintf("%s = %v (type: %v)", name, metric.GetValue(), metric.GetType()))
	}

	rw.Write([]byte(strings.Join(listOfMetrics, "\n")))
}

func (s *Server) getMetricValue(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	s.Lock()
	defer s.Unlock()
	if metric, ok := s.metrics[name]; ok {
		rw.Write([]byte(fmt.Sprintf("%v", metric.GetValue())))
		return
	}

	rw.WriteHeader(http.StatusNotFound)
}

func (s *Server) newRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
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
					r.Post("/{value}", s.postGaugeHandler)
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
					r.Post("/{value}", s.postCounterHandler)
				})
			})

			r.Route("/", func(r chi.Router) {
				r.Post("/*", postDefaultHandler)
			})
		})
	})

	return r
}
