package entity

import "math"

type Gauge struct {
	value float64
}

func (m *Gauge) GetType() string {
	return "gauge"
}

func (m *Gauge) SetValue(value any) {
	switch v := value.(type) {
	case float64:
		m.value = v
	case *float64:
		m.value = *v
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

func (m *Gauge) GetPointer() *float64 {
	return &m.value
}
