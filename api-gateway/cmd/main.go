package main

import (
	"api-gateway/internal/auth"
	"api-gateway/internal/order"
	"api-gateway/internal/product"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to Auth service
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = ":50051" // fallback for local development
	}
	log.Printf("Connecting to auth service at: %s", authServiceURL)
	authConn, err := grpc.Dial(authServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	// Connect to Product service
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = ":50052"
	}
	productConn, err := grpc.Dial(productServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Warning: Failed to connect to product service: %v", err)
	}
	defer productConn.Close()

	// Connect to Order service
	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
	if orderServiceURL == "" {
		orderServiceURL = ":50053"
	}
	orderConn, err := grpc.Dial(orderServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Warning: Failed to connect to order service: %v", err)
	}
	defer orderConn.Close()

	// Setup auth client and handler
	authClient := auth.NewAuthServiceClient(authConn)
	authHandler := auth.NewAuthHandler(authClient)

	productClient := product.NewProductServiceClient(productConn)
	productHandler := product.NewProductHandler(productClient)

	orderClient := order.NewOrderServiceClient(orderConn)
	orderHandler := order.NewOrderHandler(orderClient)
	r := gin.Default()
	middleware := auth.NewAuthMiddleware(authClient)

	// Setup auth routes
	auth.SetupAuthRoute(r, authHandler)
	product.SetupProductRoute(r, productHandler, middleware)
	order.SetupOrderRoute(r, orderHandler, middleware)
	log.Println("Port :8000 is ready to use by external services")
	r.Run(":8000")
}
