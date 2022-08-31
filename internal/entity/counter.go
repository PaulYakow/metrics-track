package entity

// todo: убрать типы counter и gauge (оставить их "под капотом" либо вообще переделать функции обработки)

type Counter struct {
	value int64
}

func (c *Counter) GetType() string {
	return "counter"
}

func (c *Counter) SetValue(value any) {
	switch v := value.(type) {
	case int64:
		c.value = v
	case *int64:
		c.value = *v
	case int:
		c.value = int64(v)
	case int32:
		c.value = int64(v)
	default:
		// todo: return error
		c.value = 0
	}
}

func (c *Counter) GetValue() int64 {
	return c.value
}

func (c *Counter) GetPointer() *int64 {
	return &c.value
}

func (c *Counter) Increment() {
	c.value++
}

func (c *Counter) IncrementDelta(delta any) {
	switch d := delta.(type) {
	case int64:
		c.value += d
	case *int64:
		c.value += *d
	case int:
		c.value += int64(d)
	case int32:
		c.value += int64(d)
	default:
		// todo: return error
		c.value = 0
	}
}
