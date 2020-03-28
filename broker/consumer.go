package broker

import (
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// MessageHandler is invoked when a message is received
type MessageHandler func(message *kafka.Message)

// Consumer represents  a message consumer
type Consumer interface {
	Close()
}

// KafkaConsumer is a Kafka message consumer
type KafkaConsumer struct {
	bootstrapServers string
	groupID          string
	topic            []string
	consumer         *kafka.Consumer
	shutdownChan     chan struct{}
	handler          MessageHandler
}

// NewKafkaConsumer creates a new Kafka Consumer and sets up message listener
func NewKafkaConsumer(bootstrapServers string, groupID string, topics []string, msgHandler MessageHandler) (*KafkaConsumer, error) {
	c := &KafkaConsumer{
		bootstrapServers: bootstrapServers,
		groupID:          groupID,
		topic:            topics,
		shutdownChan:     make(chan struct{}),
		handler:          msgHandler,
	}

	var err error
	c.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	c.consumer.SubscribeTopics(topics, nil)
	setupConsumer(c)

	return c, nil
}

func setupConsumer(c *KafkaConsumer) {
	go func() {
		for {
			select {
			case _ = <-c.shutdownChan:
				return
			default:
			}

			msg, err := c.consumer.ReadMessage(3 * time.Second)
			if msg != nil && err == nil {
				c.handler(msg)
				c.consumer.Commit()
			}
		}
	}()
}

// Close closes the consumer
func (c *KafkaConsumer) Close() {
	c.shutdownChan <- struct{}{}
	c.consumer.Close()
}
