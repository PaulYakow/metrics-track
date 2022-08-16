package repo

import (
	"errors"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"sync"
)

var errUnknownType = errors.New("unknown metric type")
var errNotFound = errors.New("not found")

// Реализация репозитория для клиента

type ClientMetricRepo struct {
	gauges   map[string]*entity.Gauge
	counters map[string]*entity.Counter
}

func NewClientRepo() *ClientMetricRepo {
	return &ClientMetricRepo{
		gauges:   make(map[string]*entity.Gauge),
		counters: make(map[string]*entity.Counter),
	}
}

func (r *ClientMetricRepo) Store(g map[string]*entity.Gauge, c map[string]*entity.Counter) {
	r.gauges = g
	r.counters = c
}

func (r *ClientMetricRepo) ReadCurrentMetrics() []string {
	result := make([]string, 0)

	for name, gauge := range r.gauges {
		result = append(result, fmt.Sprintf("/%s/%s/%v", gauge.GetType(), name, gauge.GetValue()))
	}

	for name, counter := range r.counters {
		result = append(result, fmt.Sprintf("/%s/%s/%v", counter.GetType(), name, counter.GetValue()))
	}

	return result
}

// Реализация репозитория для сервера

type ServerMetricRepo struct {
	sync.Mutex
	gauges   map[string]*entity.Gauge
	counters map[string]*entity.Counter
}

func NewServerRepo() *ServerMetricRepo {
	return &ServerMetricRepo{
		gauges:   make(map[string]*entity.Gauge),
		counters: make(map[string]*entity.Counter),
	}
}

func (r *ServerMetricRepo) Store(typeName string, name string, value any) error {
	r.Lock()
	defer r.Unlock()

	switch typeName {
	case "counter":
		if _, ok := r.counters[name]; !ok {
			r.counters[name] = &entity.Counter{}
		}
		r.counters[name].IncrementDelta(value)
		return nil
	case "gauge":
		if _, ok := r.gauges[name]; !ok {
			r.gauges[name] = &entity.Gauge{}
		}
		r.gauges[name].SetValue(value)
		return nil
	default:
		return errUnknownType
	}
}

func (r *ServerMetricRepo) ReadValue(typeName string, name string) (any, error) {
	r.Lock()
	defer r.Unlock()

	switch typeName {
	case "counter":
		if _, ok := r.counters[name]; !ok {
			return nil, errNotFound
		}
		return r.counters[name].GetValue(), nil
	case "gauge":
		if _, ok := r.gauges[name]; !ok {
			return nil, errNotFound
		}
		return r.gauges[name].GetValue(), nil
	default:
		return nil, errUnknownType
	}
}

func (r *ServerMetricRepo) ReadAll() []string {
	result := make([]string, 0)

	r.Lock()
	defer r.Unlock()

	for name, gauge := range r.gauges {
		result = append(result, fmt.Sprintf("%s = %v [%s]", name, gauge.GetValue(), gauge.GetType()))
	}

	for name, counter := range r.counters {
		result = append(result, fmt.Sprintf("%s = %v [%s]", name, counter.GetValue(), counter.GetType()))
	}

	return result
}
