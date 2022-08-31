package entity

import (
	"fmt"
	"strconv"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metric) GetValue() string {
	switch m.MType {
	case "gauge":
		return strconv.FormatFloat(*m.Value, 'f', -1, 64)
	case "counter":
		return strconv.Itoa(int(*m.Delta))
	default:
		return ""
	}
}

func (m *Metric) Update(value string) error {
	switch m.MType {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("router - update gauge value: %w", err)
		}
		m.Value = &v

	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("router - update counter value: %w", err)
		}
		*m.Delta += v

	default:
		return fmt.Errorf("router - update unknown type: %q", m.MType)
	}

	return nil
}
