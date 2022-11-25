// Package entity содержит сущность метрики и методы взаимодействия с ней.
package entity

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strconv"
)

// Metric .
type Metric struct {
	ID    string     `json:"id" db:"id"`                 // имя метрики
	MType string     `json:"type" db:"type"`             // параметр, принимающий значение gauge или counter
	Delta *int64     `json:"delta,omitempty" db:"delta"` // значение метрики в случае передачи counter
	Value *float64   `json:"value,omitempty" db:"value"` // значение метрики в случае передачи gauge
	Hash  NullString `json:"hash,omitempty" db:"hash"`   // значение хеш-функции
}

// GetValue - получение текущего значения метрики.
func (m *Metric) GetValue() []byte {
	switch m.MType {
	case "gauge":
		return strconv.AppendFloat(make([]byte, 0, 24), *m.Value, 'f', -1, 64)
	case "counter":
		return strconv.AppendInt(make([]byte, 0, 24), *m.Delta, 10)
	default:
		return []byte("")
	}
}

// Create - создание метрики по заданным параметрам.
func Create(mType, name, value string) (*Metric, error) {
	switch mType {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, &valueErr{
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
			return nil, &valueErr{
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
		return nil, &typeErr{
			name:  name,
			tName: mType,
			err:   ErrUnknownType,
		}
	}
}

// Update - обновление метрики на основе переданной.
func (m *Metric) Update(metric *Metric) error {
	switch m.MType {
	case "gauge":
		m.Value = metric.Value

	case "counter":
		*m.Delta += *metric.Delta

	default:
		return &typeErr{
			name:  m.ID,
			tName: m.MType,
			err:   ErrUnknownType,
		}
	}

	return nil
}

// UpdateValue - обновление значения метрики типа gauge.
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

// UpdateDelta - обновление значения метрики типа counter.
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

// calcHash - вычисление хэша по переданному ключу.
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

// GetHash - получить значение хэша текущей метрики.
func (m *Metric) GetHash() string {
	return string(m.Hash)
}

// SetHash - установить хэш по переданному ключу.
func (m *Metric) SetHash(key string) {
	m.Hash = NullString(m.calcHash(key))
}

// CheckHash - проверить хэш текущей метрики на соответствие переданному.
func (m *Metric) CheckHash(hash, key string) error {
	if !bytes.Equal([]byte(hash), []byte(m.calcHash(key))) {
		return &hashErr{
			name: m.ID,
			err:  ErrHashMismatch,
		}
	}

	return nil
}
