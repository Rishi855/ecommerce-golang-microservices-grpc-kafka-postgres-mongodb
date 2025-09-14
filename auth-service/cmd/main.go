package main

import (
	"log"
	"net"
	"os"

	"auth-service/internal/handler"
	"auth-service/internal/initializers"
	"auth-service/internal/proto"
	"auth-service/internal/repository"
	"auth-service/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDatabase()
}

func main() {
	port := os.Getenv("APP_PORT")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	db := initializers.ConnectDatabase()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	proto.RegisterAuthServiceServer(grpcServer, userHandler)

	log.Printf("Auth service is running on port %v...",port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
