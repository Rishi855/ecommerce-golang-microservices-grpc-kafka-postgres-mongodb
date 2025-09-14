package main

import (
	"fmt"
	"log"
	"net"

	"order-service/internal/client"
	"order-service/internal/events"
	"order-service/internal/handler"
	"order-service/internal/initializers"
	"order-service/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDatabase()
}

func main() {
	db := initializers.ConnectDatabase()
	kafkaBrokers := initializers.GetEnvWithDefault("KAFKA_BROKERS", "kafka:9092")
	notificationProducer := events.NewProducer(kafkaBrokers, events.NOTIFICATION_TOPIC)
	logProducer := events.NewProducer(kafkaBrokers, events.ORDER_LOGS_TOPIC)
	defer notificationProducer.Close()
	defer logProducer.Close()

	// connect to product service
	productServiceHost := initializers.GetEnvWithDefault("PRODUCT_SERVICE_HOST", "product-service:50052")
	productClient := client.NewProductServiceClient(productServiceHost)

	orderHandler := handler.NewOrderHandler(db, productClient, notificationProducer, logProducer)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	proto.RegisterOrderServiceServer(grpcServer, orderHandler)

	fmt.Println("Order service is running on port: 50053")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
