package broker

import (
	"os"
	"testing"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func TestConsumer(t *testing.T) {
	receivedMessage := new(*kafka.Message)
	topicName := "test-topic"
	bootstrapURL := os.Getenv("BOOTSTRAP_URLS")

	c, err := NewKafkaConsumer(bootstrapURL, "test-group", []string{topicName}, func(message *kafka.Message) {
		*receivedMessage = message
		t.Logf("Got message %s:%s", string(message.Key), string(message.Value))
	})
	if err != nil {
		t.Fatalf("Error creating consumer: %s", err)
	}
	defer c.Close()

	var p *kafka.Producer
	p, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapURL})
	if err != nil {
		t.Fatalf("Error creating producer: %s", err)
	}
	defer p.Close()

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicName},
		Key:            []byte("test"),
		Value:          []byte("test value"),
	}, nil)

	if err != nil {
		t.Fatalf("Error producing message: %s", err)
	}

	time.Sleep(10 * time.Second)

	if receivedMessage == nil {
		t.Fatalf("No message received")
	}

	k, v := string((*receivedMessage).Key), string((*receivedMessage).Value)

	if k != "test" {
		t.Fatalf("Incorrect key %s", k)
	}

	if v != "test value" {
		t.Fatalf("Incorrect value %s", v)
	}
}

func TestConsumeMultipleMessages(t *testing.T) {
	topicName := "test-topic"
	bootstrapURL := os.Getenv("BOOTSTRAP_URLS")
	msgChan := make(chan *kafka.Message, 20)

	c, err := NewKafkaConsumer(bootstrapURL, "test-group", []string{topicName}, func(message *kafka.Message) {
		msgChan <- message
		t.Logf("Got message %d:%s", message.Key[0], string(message.Value))
	})
	if err != nil {
		t.Fatalf("Error creating consumer: %s", err)
	}
	defer c.Close()

	var p *kafka.Producer
	p, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapURL})
	if err != nil {
		t.Fatalf("Error creating producer: %s", err)
	}
	defer p.Close()

	go func() {
		var i byte
		for i = 0; i < 20; i++ {
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topicName},
				Key:            []byte{i},
				Value:          []byte("test value"),
			}, nil)
		}
	}()

	time.Sleep(10 * time.Second)
	close(msgChan)

	receivedKeys := make(map[byte]struct{})
	for msg := range msgChan {
		receivedKeys[msg.Key[0]] = struct{}{}
	}

	var i byte
	for i = 0; i < 20; i++ {
		_, pres := receivedKeys[i]

		if !pres {
			t.Errorf("Key %d not present\n", i)
		}
	}
}
