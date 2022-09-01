package entity

import (
	"strconv"
)

// todo: убрать типы counter и gauge (оставить их "под капотом" либо вообще переделать функции обработки)

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

func Create(mType, name, value string) (*Metric, error) {
	switch mType {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, valueErr{
				name:  name,
				value: value,
				err:   ErrParseValue,
			}
		}
		return &Metric{
			ID:    name,
			MType: mType,
			Value: &v,
		}, nil

	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, valueErr{
				name:  name,
				value: value,
				err:   ErrParseValue,
			}
		}
		return &Metric{
			ID:    name,
			MType: mType,
			Delta: &v,
		}, nil

	default:
		return nil, typeErr{
			name:  name,
			tName: mType,
			err:   ErrUnknownType,
		}
	}
}

func (m *Metric) Update(metric *Metric) error {
	switch m.MType {
	case "gauge":
		m.Value = metric.Value

	case "counter":
		*m.Delta += *metric.Delta

	default:
		return typeErr{
			name:  m.ID,
			tName: m.MType,
			err:   ErrUnknownType,
		}
	}

	return nil
}

func (m *Metric) UpdateValue(value any) error {
	switch v := value.(type) {
	case float64:
		m.Value = &v
	case *float64:
		m.Value = v
	case uint64:
		val := float64(v)
		m.Value = &val
	case uint32:
		val := float64(v)
		m.Value = &val
	default:
		return ErrUnknownType
	}

	return nil
}

func (m *Metric) UpdateDelta(value any) error {
	switch d := value.(type) {
	case int64:
		*m.Delta += d
	case *int64:
		*m.Delta += *d
	case int:
		val := int64(d)
		*m.Delta += val
	case int32:
		val := int64(d)
		*m.Delta += val
	default:
		return ErrUnknownType
	}

	return nil
}
