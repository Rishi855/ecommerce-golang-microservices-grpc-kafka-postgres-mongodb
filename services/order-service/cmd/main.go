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
	notificationProducer := events.NewProducer("localhost:9092",events.NOTIFICATION_TOPIC)
	logProducer := events.NewProducer("localhost:9092",events.ORDER_LOGS_TOPIC)
	defer notificationProducer.Close()
	defer logProducer.Close()

	// connect to product service (running on 50052 per your setup)
	productClient := client.NewProductServiceClient("localhost:50052")

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
