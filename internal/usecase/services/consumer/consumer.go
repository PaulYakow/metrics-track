package consumer

import (
	"encoding/json"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"os"
)

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) Read() ([]entity.Metric, error) {
	defer c.file.Close()

	var metrics []entity.Metric
	if err := c.decoder.Decode(&metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}
