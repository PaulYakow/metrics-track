package model

import "math"

type Metric interface {
	GetValue() any
	SetValue(any)
	GetType() string
}

type Gauge struct {
	value float64
}

func (g *Gauge) GetValue() any {
	return g.value
}

func (g *Gauge) SetValue(v any) {
	switch i := v.(type) {
	case float64:
		g.value = i
	case uint64:
		g.value = float64(i)
	default:
		g.value = math.NaN()
	}
}

func (g *Gauge) GetType() string {
	return "gauge"
}

type Counter struct {
	value int64
}

func (c *Counter) GetValue() any {
	return c.value
}

func (c *Counter) SetValue(v any) {
	switch i := v.(type) {
	case float64:
		c.value = int64(i)
	case int64:
		c.value = i
	default:
		c.value = 0
	}
}

func (c *Counter) GetType() string {
	return "counter"
}
