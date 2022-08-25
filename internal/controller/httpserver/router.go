package httpserver

import (
	"compress/flate"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type serverRoutes struct {
	uc usecase.IServer
}

func NewRouter(uc usecase.IServer) chi.Router {
	s := &serverRoutes{uc: uc}

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Compress(flate.BestCompression))

	mux.Route("/", func(r chi.Router) {
		//Обработка GET-запроса к хосту
		r.Get("/", s.getListOfMetrics)

		r.Route("/value", func(r chi.Router) {
			r.Post("/", s.postValueByJSON)
			r.Get("/{type}/{name}", s.getMetricValue)
		})

		r.Route("/update", func(r chi.Router) {
			r.Post("/", s.postUpdateByJSON)

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
	respBody := []byte(strings.Join(s.uc.GetAllMetrics(), "\n"))
	rw.Header().Set("Content-Type", "text/html")
	rw.Write(respBody)
}

func (s *serverRoutes) getMetricValue(rw http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	value, err := s.uc.GetValueByType(mType, name)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	respBody := []byte(value)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respBody)
}

func (s *serverRoutes) postValueByJSON(rw http.ResponseWriter, r *http.Request) {
	// 1. В теле запроса JSON с ID и MType
	// 2. Заполнить значение метрики
	// 3. Отправить ответный JSON
	if r.Header.Get("Content-Type") != "application/json" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	respBody, err := s.uc.GetValueByJSON(reqBody)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respBody)
}

func (s *serverRoutes) postUpdateByJSON(rw http.ResponseWriter, r *http.Request) {
	// Обработать JSON из тела запроса - сохранить в соответствующую метрику переданное значение
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	if r.Header.Get("Content-Type") != "application/json" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.uc.SaveValueByJSON(body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
