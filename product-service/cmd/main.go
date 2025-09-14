package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"product-service/internal/handler"
	"product-service/internal/initializers"
	"product-service/internal/proto"
	"product-service/internal/repository"
	"product-service/internal/service"

	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDatabase() // optional warm-up
}

func main() {
	port := os.Getenv("APP_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db := initializers.ConnectDatabase()
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	proto.RegisterProductServiceServer(grpcServer, h)

	fmt.Println("Product service is running on :"+port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
