package global

import (
	"errors"
	"github.com/crask/mqproxy/producer/kafka"
	"sync"
)

var ErrInvalidParam error = errors.New("Invalid offset when close producer in pool")

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

func (pool *KafkaProducerPool) GetProducer() (*producer.KafkaProducer, int) {
	pool.curr++
	pos := pool.curr % pool.size
	pool.curr = pos
	return pool.producers[pos], pos
}

func (pool *KafkaProducerPool) ReopenProducer(kp *producer.KafkaProducer, i int) error {
	if i >= pool.size {
		return ErrInvalidParam
	}

	pool.lock.Lock()
	defer pool.lock.Unlock()

	pool.producers[i].Close()
	//pool.producers[i] = nil
	var err error
	if pool.producers[i], err = producer.NewKafkaProducer(pool.config); err != nil {
		return err
	}
	return nil
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
