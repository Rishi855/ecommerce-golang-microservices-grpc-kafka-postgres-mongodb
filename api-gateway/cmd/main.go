package main

import (
	"api-gateway/internal/auth"
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	// Connect to Auth service
	authConn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	// Connect to Product service
	productConn, err := grpc.Dial(":50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}
	defer productConn.Close()

	// Connect to Order service
	// orderConn, err := grpc.Dial(":50053", grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("Failed to connect to order service: %v", err)
	// }
	// defer orderConn.Close()

	// Setup auth client and handler
	authClient := auth.NewAuthServiceClient(authConn)
	authHandler := auth.NewAuthHandler(authClient)

	r := gin.Default()
	// middleware := auth.NewAuthMiddleware(authClient)

	// Setup auth routes
	auth.SetupAuthRoute(r, authHandler)
	log.Println("Port :8000 is ready to use by external services")
	r.Run(":8000")
}
