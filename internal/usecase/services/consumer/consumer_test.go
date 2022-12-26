package consumer

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConsumer(t *testing.T) {
	var buf bytes.Buffer

	// Корректные данные
	buf.WriteString(`[
{
	"id": "test",
	"type": "testT",
	"value": 99.9
}
]`)
	data, err := readFile(&buf)
	require.NoError(t, err)
	require.NotNil(t, data)

	c := &Consumer{decoder: json.NewDecoder(bytes.NewReader(data))}
	metrics, err := c.Read()
	require.NoError(t, err)
	require.NotNil(t, metrics)

	// Некорректные данные
	buf.Reset()
	buf.WriteString(`
{
	"id": "test",
	"type": "testT",
	"value": 99.9
}`)
	data, err = readFile(&buf)
	require.NoError(t, err)
	require.NotNil(t, data)

	c = &Consumer{decoder: json.NewDecoder(bytes.NewReader(data))}
	metrics, err = c.Read()
	require.Error(t, err)
	require.Empty(t, metrics)

	// Некорректный путь до файла
	c, err = NewConsumer("/wrong/path/to/file.db")
	require.Error(t, err)
	require.Empty(t, c)
}
