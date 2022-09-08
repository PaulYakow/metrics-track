package entity

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strconv"
)

// todo: убрать типы counter и gauge (оставить их "под капотом" либо вообще переделать функции обработки)

type Metric struct {
	ID    string     `json:"id" db:"id"`                 // имя метрики
	MType string     `json:"type" db:"type"`             // параметр, принимающий значение gauge или counter
	Delta *int64     `json:"delta,omitempty" db:"delta"` // значение метрики в случае передачи counter
	Value *float64   `json:"value,omitempty" db:"value"` // значение метрики в случае передачи gauge
	Hash  NullString `json:"hash,omitempty" db:"hash"`   // значение хеш-функции
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

func (m *Metric) UpdateValue(value any) {
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
		m.Value = nil
	}
}

func (m *Metric) UpdateDelta(value any) {
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
		m.Delta = nil
	}
}

func (m *Metric) calcHash(key string) string {
	var data string
	switch m.MType {
	case "counter":
		data = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
	case "gauge":
		data = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (m *Metric) GetHash() string {
	return string(m.Hash)
}

func (m *Metric) SetHash(key string) {
	m.Hash = NullString(m.calcHash(key))
}

func (m *Metric) CheckHash(hash, key string) error {
	if !bytes.Equal([]byte(hash), []byte(m.calcHash(key))) {
		return ErrHashMismatch
	}

	return nil
}
