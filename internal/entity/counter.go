package entity

// todo: убрать типы counter и gauge (оставить их "под капотом" либо вообще переделать функции обработки)

type Counter struct {
	value int64
}

func (m *Counter) GetType() string {
	return "counter"
}

func (m *Counter) SetValue(value any) {
	switch v := value.(type) {
	case int64:
		m.value = v
	case *int64:
		m.value = *v
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

func (m *Counter) GetPointer() *int64 {
	return &m.value
}

func (m *Counter) Increment() {
	m.value++
}

func (m *Counter) IncrementDelta(delta any) {
	switch d := delta.(type) {
	case int64:
		m.value += d
	case *int64:
		m.value += *d
	case int:
		m.value += int64(d)
	case int32:
		m.value += int64(d)
	default:
		m.value = 0
	}
}
