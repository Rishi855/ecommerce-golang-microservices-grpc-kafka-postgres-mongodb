package events

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct{
	reader *kafka.Reader
}

func NewConsumer(broker, topic, groupId string) *Consumer{
	return &Consumer{
		reader : kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic: topic,
			GroupID: groupId,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}
}

func (c *Consumer) Consume(ctx context.Context, handler func(key, value string) error){
	for{
		log.Println("\n#####################################################################################################")
		msg, err:= c.reader.ReadMessage(ctx)
		if err!=nil{
			log.Printf("Error while reading message: %v",err)
			continue
		}
		log.Printf("Received message key=%s value=%s",string(msg.Key),string(msg.Value))

		if err:=handler(string(msg.Key),string(msg.Value)); err!=nil{
			log.Printf("Error handling message: %v",err)
		}
	}
}

func (c *Consumer) Close() error{
	return c.reader.Close()
}