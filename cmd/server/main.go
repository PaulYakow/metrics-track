package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"strconv"
)

type gauge float64
type counter int64

var gaugeMetrics = make(map[string]gauge)
var counterMetrics = make(map[string]counter)

func gaugeHandler(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	rawValue := chi.URLParam(r, "value")

	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	gaugeMetrics[name] = gauge(value)

	rw.WriteHeader(http.StatusOK)
}

func counterHandler(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	rawValue := chi.URLParam(r, "value")

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	counterMetrics[name] += counter(value)

	rw.WriteHeader(http.StatusOK)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func getListOfMetrics(rw http.ResponseWriter, r *http.Request) {
	gauges, _ := json.Marshal(gaugeMetrics)

	rw.Write(gauges)
}

func getMetricValue(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	switch metricType {
	case "gauge":
		if value, ok := gaugeMetrics[name]; ok {
			rw.Write([]byte(fmt.Sprintf("%v", value)))
			return
		}
		rw.WriteHeader(http.StatusNotFound)
	case "counter":
		if value, ok := counterMetrics[name]; ok {
			rw.Write([]byte(fmt.Sprintf("%v", value)))
			return
		}
		rw.WriteHeader(http.StatusNotFound)
	default:
		rw.WriteHeader(http.StatusNotFound)
	}
}

func newRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		//Обработка GET-запроса к хосту
		r.Get("/", getListOfMetrics)

		r.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}", getMetricValue)
		})

		r.Route("/update", func(r chi.Router) {

			r.Route("/gauge/{name}", func(r chi.Router) {
				r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
					rw.WriteHeader(http.StatusBadRequest)
				})
				r.Post("/{value}", gaugeHandler)
			})

			r.Route("/counter/{name}", func(r chi.Router) {
				r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
					rw.WriteHeader(http.StatusBadRequest)
				})
				r.Post("/{value}", counterHandler)
			})

			r.Route("/", func(r chi.Router) {
				r.Post("/*", defaultHandler)
			})
		})
	})

	return r
}

func main() {
	router := newRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
