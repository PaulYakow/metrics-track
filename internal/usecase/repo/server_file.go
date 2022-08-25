package repo

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/consumer"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/producer"
	"log"
	"sync"
)

type ServerFile struct {
	sync.Mutex
	producer *producer.Producer
	consumer *consumer.Consumer
}

func NewServerFile(filename string) (*ServerFile, error) {
	p, err := producer.NewProducer(filename)
	if err != nil {
		return nil, err
	}

	c, err := consumer.NewConsumer(filename)
	if err != nil {
		return nil, err
	}

	return &ServerFile{
		producer: p,
		consumer: c,
	}, nil
}

func (repo *ServerFile) SaveMetrics(metrics []entity.Metrics) {
	repo.Lock()
	defer repo.Unlock()

	err := repo.producer.WriteMetric(&metrics)
	if err != nil {
		log.Println(err)
	}
}

func (repo *ServerFile) ReadMetrics() []entity.Metrics {
	repo.Lock()
	defer repo.Unlock()

	metrics, err := repo.consumer.ReadMetric()
	if err != nil {
		log.Println(err)
	}

	return metrics
}
