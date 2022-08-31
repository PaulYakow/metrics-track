package entity

import "math"

type Gauge struct {
	value float64
}

func (g *Gauge) GetType() string {
	return "gauge"
}

func (g *Gauge) SetValue(value any) {
	switch v := value.(type) {
	case float64:
		g.value = v
	case *float64:
		g.value = *v
	case uint64:
		g.value = float64(v)
	case uint32:
		g.value = float64(v)
	default:
		// todo: return error
		g.value = math.NaN()
	}
}

func (g *Gauge) GetValue() float64 {
	return g.value
}

func (g *Gauge) GetPointer() *float64 {
	return &g.value
}
