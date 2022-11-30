// Package consumer сервис чтения метрик из файлов JSON.
package consumer

import (
	"encoding/json"
	"os"

	"github.com/PaulYakow/metrics-track/internal/entity"
)

// Consumer читатель данных (JSON).
type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

// NewConsumer создаёт объект Consumer.
// В качестве параметра принимает путь к файлу.
func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// Read - возвращает считанный из файла массив.
func (c *Consumer) Read() ([]*entity.Metric, error) {
	defer c.file.Close()

	var metrics []*entity.Metric
	if err := c.decoder.Decode(&metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}
