// Package entity определяет основные объекты для бизнес-логики (сервисы), отображение базы данных и объекты ответа HTTP, если это возможно.
// Каждая логическая группа объектов в собственном файле.
package entity

import "math"

type Gauge struct {
	value float64
}

type Counter struct {
	value int64
}

// Методы для типа Gauge

func (m *Gauge) GetType() string {
	return "gauge"
}

func (m *Gauge) SetValue(value any) {
	switch v := value.(type) {
	case float64:
		m.value = v
	case uint64:
		m.value = float64(v)
	case uint32:
		m.value = float64(v)
	default:
		m.value = math.NaN()
	}
}

func (m *Gauge) GetValue() float64 {
	return m.value
}

// Методы для типа Counter

func (m *Counter) GetType() string {
	return "counter"
}

func (m *Counter) SetValue(value any) {
	switch v := value.(type) {
	case int:
		m.value = int64(v)
	case int32:
		m.value = int64(v)
	default:
		m.value = 0
	}
}

func (m *Counter) GetValue() int64 {
	return m.value
}

func (m *Counter) Increment() {
	m.value++
}

func (m *Counter) IncrementDelta(delta int64) {
	m.value += delta
}
