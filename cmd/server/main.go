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
	switch r.Method {
	case http.MethodPost:
		metricFromPath := strings.Split(strings.TrimLeft(r.URL.Path, "/updategauge"), "/")

		value, _ := strconv.ParseFloat(metricFromPath[1], 64)
		metrics[metricFromPath[0]] = metric{
			tMetric: "gauge",
			value:   gauge(value),
		}
	}
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		metricFromPath := strings.Split(strings.TrimLeft(r.URL.Path, "/updatecounter"), "/")

		value, _ := strconv.Atoi(metricFromPath[1])
		metrics[metricFromPath[0]] = metric{
			tMetric: "gauge",
			value:   counter(value),
		}
	}
}

func main() {
	http.HandleFunc("/update/gauge/", gaugeHandler)
	http.HandleFunc("/update/counter/", counterHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
