package events

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(broker),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
			Async:        false,
		},
	}
}

func (p *Producer) Publish(key, value string) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
		Time:  time.Now(),
	}
	err := p.writer.WriteMessages(context.Background(),msg)
	if err!=nil{
		log.Printf("failed to publish message: %v",err)
		return err
	}
	log.Printf("publish message to kafka %s",value)
	return nil
}

func (p *Producer) Close() error{
	return p.writer.Close()
}
