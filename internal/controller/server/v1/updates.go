package v1

import (
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"strconv"
)

type updates struct{}

func (rs updates) Routes(s *serverRoutes) chi.Router {
	r := chi.NewRouter()

	r.Post("/", s.updateByJSON)

	r.Post("/*", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotImplemented)
	})

	r.Route("/gauge", func(r chi.Router) {
		r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
		})
		r.Post("/{name}/", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusBadRequest)
		})
		r.Post("/{name}/{value}", s.updateGaugeByURL)
	})

	r.Route("/counter", func(r chi.Router) {
		r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
		})
		r.Post("/{name}/", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusBadRequest)
		})
		r.Post("/{name}/{value}", s.updateCounterByURL)
	})

	return r
}

func (s *serverRoutes) updateGaugeByURL(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	rawValue := chi.URLParam(r, "value")

	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		log.Printf("post gauge value: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	s.uc.SaveGauge(name, value)
	rw.WriteHeader(http.StatusOK)
}

func (s *serverRoutes) updateCounterByURL(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	rawValue := chi.URLParam(r, "value")

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		log.Printf("post counter value: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	s.uc.SaveCounter(name, value)
	rw.WriteHeader(http.StatusOK)
}

func (s *serverRoutes) updateByJSON(rw http.ResponseWriter, r *http.Request) {
	// Обработать JSON из тела запроса - сохранить в соответствующую метрику переданное значение
	if r.Header.Get("Content-Type") != "application/json" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("read request body %q: %v", r.URL.Path, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.uc.SaveValueByJSON(body)
	if err != nil {
		log.Printf("save value to storage: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
