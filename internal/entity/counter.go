package entity

// todo: убрать типы counter и gauge (оставить их "под капотом" либо вообще переделать функции обработки)

type Counter struct {
	delta int64
}

func (c *Counter) GetType() string {
	return "counter"
}

func (c *Counter) SetValue(value any) {
	switch v := value.(type) {
	case int64:
		c.delta = v
	case *int64:
		c.delta = *v
	case int:
		c.delta = int64(v)
	case int32:
		c.delta = int64(v)
	default:
		// todo: return error
		c.delta = 0
	}
}

func (c *Counter) GetValue() int64 {
	return c.delta
}

func (c *Counter) GetPointer() *int64 {
	return &c.delta
}

func (c *Counter) Increment() {
	c.delta++
}

func (c *Counter) IncrementDelta(delta any) {
	switch d := delta.(type) {
	case int64:
		c.delta += d
	case *int64:
		c.delta += *d
	case int:
		c.delta += int64(d)
	case int32:
		c.delta += int64(d)
	default:
		// todo: return error
		c.delta = 0
	}
}
