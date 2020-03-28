package broker

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Producer represents a message producer
type Producer interface {
	Produce(key []byte, message []byte)
	Close()
}

// KafkaProducer is a Kafka message producer
type KafkaProducer struct {
	bootstrapServers []string
	producer         *kafka.Producer
}

// NewKafkaProducer returns a new KafkaProducer
func NewKafkaProducer(bootstrapServers []string) (*KafkaProducer, error) {
	p := &KafkaProducer{
		bootstrapServers: bootstrapServers,
	}

	var err error
	p.producer, err = kafka.NewProducer(&kafka.ConfigMap{"boostrap.servers": bootstrapServers})
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Produce emits a message
func (p *KafkaProducer) Produce(key []byte, value []byte) {
	p.producer.Produce(&kafka.Message{
		Key:   key,
		Value: value,
	}, nil)
}
