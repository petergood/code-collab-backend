package broker

import (
	"os"
	"testing"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func TestProducer(t *testing.T) {
	bootstrapURL := os.Getenv("BOOTSTRAP_URLS")
	p, err := NewKafkaProducer(bootstrapURL, "test-topic")
	if err != nil {
		t.Errorf("Cannot create producer: %s", err)
		return
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapURL,
		"group.id":          "test-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		t.Errorf("Could not create consumer: %s", err)
		return
	}

	c.SubscribeTopics([]string{"test-topic"}, nil)

	stop := make(chan struct{}, 1)
	received := false
	go func() {
		for {
			select {
			case _ = <-stop:
				t.Errorf("Consumer timed out")
				break
			default:
			}

			msg, err := c.ReadMessage(3 * time.Second)
			if err == nil {
				k, v := string(msg.Key), string(msg.Value)

				if k != "asd" {
					t.Errorf("Incorrect key %s", k)
				} else if v != "Hello world!" {
					t.Errorf("Incorrect value %v", v)
				} else {
					received = true
				}

				break
			}
		}
	}()

	err = p.Produce([]byte("asd"), []byte("Hello world!"))

	if err != nil {
		t.Errorf("Error producing message: %s", err)
	}

	time.Sleep(10 * time.Second)
	if !received {
		stop <- struct{}{}
	}
}
