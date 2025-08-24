package main

import (
	"fmt"
	"log"
	"net"

	"order-service/internal/client"
	"order-service/internal/handler"
	"order-service/internal/initializers"
	"order-service/internal/proto"

	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDatabase()
}

func main() {
	db := initializers.ConnectDatabase()

	// connect to product service (running on 50052 per your setup)
	productClient := client.NewProductServiceClient("localhost:50052")

	orderHandler := handler.NewOrderHandler(db, productClient)

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
