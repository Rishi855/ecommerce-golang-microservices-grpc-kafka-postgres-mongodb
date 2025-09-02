package main

import (
	"context"
	"log"
	"notification-service/internal/events"
	"notification-service/internal/handler"
	"notification-service/internal/initializers"
	"os"
	"os/signal"
	"syscall"
)

func init(){
	initializers.LoadEnv()
}

func main(){
	broker := initializers.GetEnv("KAFKA_BROKER", "localhost:9092")
	topic := initializers.GetEnv("KAFKA_TOPIC", "notification-events")
	groupID := initializers.GetEnv("KAFKA_GROUP", "notification-service-group")

	consumer := events.NewConsumer(broker,topic,groupID)
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())

	go func ()  {
		consumer.Consume(ctx,handler.HandleNotification)	
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down notification service...")
	cancel()
}