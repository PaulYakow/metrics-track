// Package consumer сервис чтения метрик из файлов JSON.
package consumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/PaulYakow/metrics-track/internal/entity"
)

// Consumer читатель данных (JSON).
type Consumer struct {
	//file    *os.File
	decoder *json.Decoder
}

// NewConsumer создаёт объект Consumer.
// В качестве параметра принимает путь к файлу.
func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	data, err := readFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	defer file.Close()

	return &Consumer{
		//file:    file,
		decoder: json.NewDecoder(bytes.NewReader(data)),
	}, nil
}

// Read - возвращает считанный из файла массив.
func (c *Consumer) Read() ([]*entity.Metric, error) {
	//defer c.file.Close()

	var metrics []*entity.Metric
	if err := c.decoder.Decode(&metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}

func readFile(reader io.Reader) ([]byte, error) {
	return io.ReadAll(reader)
}
