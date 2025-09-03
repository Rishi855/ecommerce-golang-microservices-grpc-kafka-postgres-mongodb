package events

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(broker, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{broker},
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3,  // 10KB
			MaxBytes: 10e6,  // 10MB
		}),
	}
}

func (c *Consumer) Consume(ctx context.Context, handler func(key, value string) error) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("[Log-Service] Error while reading message: %v", err)
			continue
		}
		log.Printf("[Log-Service] Received log event: key=%s value=%s", string(msg.Key), string(msg.Value))

		if err := handler(string(msg.Key), string(msg.Value)); err != nil {
			log.Printf("[Log-Service] Error handling log: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
