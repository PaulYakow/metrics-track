package repo

import (
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/consumer"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/producer"
	"log"
	"sync"
)

type serverFile struct {
	sync.Mutex
	producer *producer.Producer
	consumer *consumer.Consumer
}

func NewServerFile(filename string) (*serverFile, error) {
	p, err := producer.NewProducer(filename)
	if err != nil {
		return nil, err
	}

	c, err := consumer.NewConsumer(filename)
	if err != nil {
		return nil, err
	}

	return &serverFile{
		producer: p,
		consumer: c,
	}, nil
}

func (repo *serverFile) SaveMetrics(metrics []entity.Metric) {
	repo.Lock()
	defer repo.Unlock()

	err := repo.producer.Write(&metrics)
	if err != nil {
		log.Printf("save in file: %v", err)
	}
}

func (repo *serverFile) ReadMetrics() []*entity.Metric {
	repo.Lock()
	defer repo.Unlock()

	metrics, err := repo.consumer.Read()
	if err != nil {
		log.Printf("read from file: %v", err)
	}

	return metrics
}

func (repo *serverFile) CheckConnection() error {
	return fmt.Errorf("not implement to file storage")
}
