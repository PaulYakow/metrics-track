package entity

import (
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

func (m *Metric) Update(in *Metric) error {
	switch m.MType {
	case "gauge":
		m.Value = in.Value

	case "counter":
		*m.Delta += *in.Delta

	default:
		return typeErr{
			name:  m.ID,
			tName: m.MType,
			err:   ErrUnknownType,
		}
	}

	return nil
}
