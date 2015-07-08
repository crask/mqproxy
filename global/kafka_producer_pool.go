package global

import (
	"github.com/crask/mqproxy/producer/kafka"
	"sync"
)

type KafkaProducerPool struct {
	config    *producer.KafkaProducerConfig
	producers []*producer.KafkaProducer
	size      int
	curr      int
	lock      sync.Mutex
}

func NewKafkaProducerPool(config *producer.KafkaProducerConfig, poolSize int) (*KafkaProducerPool, error) {
	var err error

	pool := &KafkaProducerPool{
		config:    config,
		producers: make([]*producer.KafkaProducer, poolSize),
		size:      0,
		curr:      -1,
	}

	for i := 0; i < poolSize; i++ {
		if pool.producers[i], err = producer.NewKafkaProducer(config); err != nil {
			return pool, err
		}
		pool.size++
	}
	ProducerPool = pool

	return pool, nil
}

func DestoryKafkaProducerPool(pool *KafkaProducerPool) error {
	for i := 0; i < pool.size; i++ {
		pool.producers[i].Close()
	}

	return nil
}

func (pool *KafkaProducerPool) GetProducer() *producer.KafkaProducer {
	return pool.producers[(pool.curr+1)%pool.size]
}

func (pool *KafkaProducerPool) Rebuild() error {
	var err error

	pool.lock.Lock()
	defer pool.lock.Unlock()

	for i := 0; i < pool.size; i++ {
		pool.producers[i].Close()
	}
	for i := 0; i < pool.size; i++ {
		if pool.producers[i], err = producer.NewKafkaProducer(pool.config); err != nil {
			return err
		}
	}

	return nil
}
