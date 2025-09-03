package main

import (
	"context"
	"log"
	"log-service/internal/events"
	"log-service/internal/handler"
	"log-service/internal/initializers"
	"os"
	"os/signal"
	"syscall"
)
func init(){
	initializers.LoadEnv()
}

func main() {

	broker := initializers.GetEnv("KAFKA_BROKER", "localhost:9092")
	topic := initializers.GetEnv("KAFKA_TOPIC", "logs.order-service")
	groupID := initializers.GetEnv("KAFKA_GROUP", "log-service-group")
	mongoURI := initializers.GetEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := initializers.GetEnv("MONGO_DB", "logsdb")
	mongoCollection := initializers.GetEnv("MONGO_COLLECTION", "order_logs")

	initializers.ConnectMongoDB(mongoURI, mongoDB, mongoCollection)

	consumer := events.NewConsumer(broker, topic, groupID)
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		consumer.Consume(ctx, handler.HandleLog)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("[Log-Service] Shutting down gracefully...")
	cancel()
}
