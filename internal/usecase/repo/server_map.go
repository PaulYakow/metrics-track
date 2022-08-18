package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"log"
	"math"
	"sync"
)

var errUnknownType = errors.New("unknown metric type")
var errNotFound = errors.New("not found")

type ServerRepo struct {
	sync.Mutex
	gauges   map[string]*entity.Gauge
	counters map[string]*entity.Counter
}

func NewServerRepo() *ServerRepo {
	return &ServerRepo{
		gauges:   make(map[string]*entity.Gauge),
		counters: make(map[string]*entity.Counter),
	}
}

func (r *ServerRepo) Store(mType string, name string, value any) error {
	r.Lock()
	defer r.Unlock()

	switch mType {
	case "gauge":
		if _, ok := r.gauges[name]; !ok {
			r.gauges[name] = &entity.Gauge{}
		}
		r.gauges[name].SetValue(value)
		return nil

	case "counter":
		if _, ok := r.counters[name]; !ok {
			r.counters[name] = &entity.Counter{}
		}
		r.counters[name].IncrementDelta(value)
		return nil
	default:
		return errUnknownType
	}
}

func (r *ServerRepo) StoreByJson(data []byte) error {
	metric := entity.Metrics{}
	if err := json.Unmarshal(data, &metric); err != nil {
		log.Fatal(err)
	}

	switch metric.MType {
	case "gauge":
		if _, ok := r.gauges[metric.ID]; !ok {
			r.gauges[metric.ID] = &entity.Gauge{}
		}
		r.gauges[metric.ID].SetValue(metric.Value)
		return nil

	case "counter":
		if _, ok := r.counters[metric.ID]; !ok {
			r.counters[metric.ID] = &entity.Counter{}
		}
		r.counters[metric.ID].IncrementDelta(metric.Delta)
		return nil
	default:
		return errUnknownType
	}
}

func (r *ServerRepo) ReadValueByType(typeName string, name string) (any, error) {
	r.Lock()
	defer r.Unlock()

	switch typeName {
	case "counter":
		return r.readCounter(name)
	case "gauge":
		return r.readGauge(name)
	default:
		return nil, errUnknownType
	}
}

func (r *ServerRepo) ReadValueByJson(data []byte) ([]byte, error) {
	metric := entity.Metrics{}
	if err := json.Unmarshal(data, &metric); err != nil {
		log.Fatal(err)
	}

	switch metric.MType {
	case "gauge":
		value, err := r.readGauge(metric.ID)
		if err != nil {
			return nil, err
		}
		metric.Value = &value

	case "counter":
		value, err := r.readCounter(metric.ID)
		if err != nil {
			return nil, err
		}

		metric.Delta = &value

	default:
		return nil, errUnknownType
	}

	result, err := json.Marshal(metric)
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}

func (r *ServerRepo) readGauge(name string) (float64, error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.gauges[name]; !ok {
		return math.NaN(), errNotFound
	}
	return r.gauges[name].GetValue(), nil
}

func (r *ServerRepo) readCounter(name string) (int64, error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.counters[name]; !ok {
		return 0, errNotFound
	}
	return r.counters[name].GetValue(), nil
}

func (r *ServerRepo) ReadAll() []string {
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
