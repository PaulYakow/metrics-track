package repo

import (
	"encoding/json"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"log"
	"sync"
)

type ClientRepo struct {
	sync.Mutex
	gauges   map[string]*entity.Gauge
	counters map[string]*entity.Counter
}

func NewClientRepo() *ClientRepo {
	return &ClientRepo{
		gauges:   make(map[string]*entity.Gauge),
		counters: make(map[string]*entity.Counter),
	}
}

func (r *ClientRepo) Store(g map[string]*entity.Gauge, c map[string]*entity.Counter) {
	r.Lock()
	defer r.Unlock()

	r.gauges = g
	r.counters = c
}

func (r *ClientRepo) ReadCurrentMetrics() []string {
	result := make([]string, 0)

	r.Lock()
	defer r.Unlock()

	for name, gauge := range r.gauges {
		result = append(result, fmt.Sprintf("/%s/%s/%v", gauge.GetType(), name, gauge.GetValue()))
	}

	for name, counter := range r.counters {
		result = append(result, fmt.Sprintf("/%s/%s/%v", counter.GetType(), name, counter.GetValue()))
	}

	return result
}

func (r *ClientRepo) ReadCurrentValues() [][]byte {
	result := make([][]byte, 0)
	metric := entity.Metric{}

	r.Lock()
	defer r.Unlock()

	for name, gauge := range r.gauges {
		metric.ID = name
		metric.MType = "gauge"
		metric.Value = gauge.GetPointer()

		data, err := json.Marshal(metric)
		if err != nil {
			log.Printf("read gauge: %v", err)
		}
		result = append(result, data)
	}

	for name, counter := range r.counters {
		metric.ID = name
		metric.MType = "counter"
		metric.Delta = counter.GetPointer()

		data, err := json.Marshal(metric)
		if err != nil {
			log.Printf("read counter: %v", err)
		}
		result = append(result, data)
	}

	return result
}
