package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type gauge float64
type counter int64

type metric struct {
	tMetric string
	value   any
}

var metrics = make(map[string]metric)

func gaugeHandler(w http.ResponseWriter, r *http.Request) {
	metricFromPath := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	switch {
	case len(metricFromPath) == 1:
		w.WriteHeader(http.StatusNotFound)
	case len(metricFromPath) == 2:
		w.WriteHeader(http.StatusNotFound)
	case len(metricFromPath) == 3:
		w.WriteHeader(http.StatusBadRequest)
	case len(metricFromPath) == 4:
		switch r.Method {
		case http.MethodPost:
			value, err := strconv.ParseFloat(metricFromPath[3], 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			metrics[metricFromPath[2]] = metric{
				tMetric: "gauge",
				value:   gauge(value),
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	metricFromPath := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	switch {
	case len(metricFromPath) == 1:
		w.WriteHeader(http.StatusNotFound)
	case len(metricFromPath) == 2:
		w.WriteHeader(http.StatusNotFound)
	case len(metricFromPath) == 3:
		w.WriteHeader(http.StatusBadRequest)
	case len(metricFromPath) == 4:
		switch r.Method {
		case http.MethodPost:
			value, err := strconv.Atoi(metricFromPath[3])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			metrics[metricFromPath[2]] = metric{
				tMetric: "gauge",
				value:   counter(value),
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func main() {
	http.HandleFunc("/update/gauge/", gaugeHandler)
	http.HandleFunc("/update/counter/", counterHandler)
	http.HandleFunc("/update/", defaultHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
