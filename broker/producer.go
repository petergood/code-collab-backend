package broker

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Producer represents a message producer
type Producer interface {
	Produce(key []byte, message []byte) error
	Close()
}

// KafkaProducer is a Kafka message producer
type KafkaProducer struct {
	bootstrapServers string
	topicName        string
	producer         *kafka.Producer
}

// NewKafkaProducer returns a new KafkaProducer
func NewKafkaProducer(bootstrapServers string, topicName string) (*KafkaProducer, error) {
	p := &KafkaProducer{
		bootstrapServers: bootstrapServers,
		topicName:        topicName,
	}

	var err error
	p.producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Produce emits a message
func (p *KafkaProducer) Produce(key []byte, value []byte) error {
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topicName, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}, nil)

	if err != nil {
		return err
	}

	return nil
}

// Close closes the producer
func (p *KafkaProducer) Close() {
	p.Close()
}
