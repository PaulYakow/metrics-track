package repo

import (
	"encoding/json"
	"errors"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"log"
	"math"
	"sync"
)

var errUnknownType = errors.New("unknown metric type")
var errNotFound = errors.New("not found")

type ServerMemory struct {
	sync.Mutex
	gauges   map[string]*entity.Gauge
	counters map[string]*entity.Counter
}

func NewServerMemory() *ServerMemory {
	return &ServerMemory{
		gauges:   make(map[string]*entity.Gauge),
		counters: make(map[string]*entity.Counter),
	}
}

func (repo *ServerMemory) InitializeMetrics(metrics []entity.Metrics) {
	var value any
	for _, metric := range metrics {
		switch metric.MType {
		case "gauge":
			value = metric.Value
		case "counter":
			value = metric.Delta
		}

		err := repo.Store(metric.MType, metric.ID, value)
		if err != nil {
			log.Printf("init metrics: %v", err)
		}
	}
}

func (repo *ServerMemory) Store(mType string, name string, value any) error {
	switch mType {
	case "gauge":
		repo.storeGauge(name, value)
		return nil

	case "counter":
		repo.storeCounter(name, value)
		return nil
	default:
		return errUnknownType
	}
}

func (repo *ServerMemory) StoreByJSON(data []byte) error {
	metric := entity.Metrics{}
	if err := json.Unmarshal(data, &metric); err != nil {
		log.Printf("store by json: %v", err)
		return err
	}

	switch metric.MType {
	case "gauge":
		repo.storeGauge(metric.ID, metric.Value)
		return nil

	case "counter":
		repo.storeCounter(metric.ID, metric.Delta)
		return nil
	default:
		return errUnknownType
	}
}

func (repo *ServerMemory) ReadValueByType(mType string, name string) (any, error) {
	switch mType {
	case "counter":
		return repo.readCounter(name)
	case "gauge":
		return repo.readGauge(name)
	default:
		return nil, errUnknownType
	}
}

func (repo *ServerMemory) ReadValueByJSON(data []byte) ([]byte, error) {
	metric := entity.Metrics{}
	if err := json.Unmarshal(data, &metric); err != nil {
		log.Printf("read value by json (unmarshal): %v", err)
		return nil, err
	}

	switch metric.MType {
	case "gauge":
		value, err := repo.readGauge(metric.ID)
		if err != nil {
			log.Printf("read gauge %q: %v", metric.ID, err)
			return nil, err
		}
		metric.Value = &value

	case "counter":
		value, err := repo.readCounter(metric.ID)
		if err != nil {
			log.Printf("read counter %q: %v", metric.ID, err)
			return nil, err
		}

		metric.Delta = &value

	default:
		return nil, errUnknownType
	}

	result, err := json.Marshal(metric)
	if err != nil {
		log.Printf("read value by json (marshal): %v", err)
		return nil, err
	}

	return result, nil
}

func (repo *ServerMemory) ReadAll() (map[string]*entity.Gauge, map[string]*entity.Counter) {
	repo.Lock()
	defer repo.Unlock()

	return repo.gauges, repo.counters
}

func (repo *ServerMemory) readGauge(name string) (float64, error) {
	repo.Lock()
	defer repo.Unlock()

	if _, ok := repo.gauges[name]; !ok {
		return math.NaN(), errNotFound
	}
	return repo.gauges[name].GetValue(), nil
}

func (repo *ServerMemory) readCounter(name string) (int64, error) {
	repo.Lock()
	defer repo.Unlock()

	if _, ok := repo.counters[name]; !ok {
		return 0, errNotFound
	}
	return repo.counters[name].GetValue(), nil
}

func (repo *ServerMemory) storeGauge(name string, value any) {
	repo.Lock()
	defer repo.Unlock()

	if _, ok := repo.gauges[name]; !ok {
		repo.gauges[name] = &entity.Gauge{}
	}
	repo.gauges[name].SetValue(value)
}

func (repo *ServerMemory) storeCounter(name string, value any) {
	repo.Lock()
	defer repo.Unlock()

	if _, ok := repo.counters[name]; !ok {
		repo.counters[name] = &entity.Counter{}
	}
	repo.counters[name].IncrementDelta(value)
}
