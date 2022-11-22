package producer

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/PaulYakow/metrics-track/internal/entity"
)

// Producer писатель данных (JSON).
type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

// NewProducer создаёт объект Producer.
// В качестве параметра принимает путь к файлу.
func NewProducer(filename string) (*Producer, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(filename), 0644)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0644)
	if err != nil {
		return nil, err
	}
	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// Write - записывает переданный массив в файл.
func (p *Producer) Write(metrics *[]entity.Metric) error {
	defer p.file.Close()

	return p.encoder.Encode(&metrics)
}
