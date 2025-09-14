package client

import (
	"context"
	"log"
	"time"

	"order-service/internal/proto"

	"google.golang.org/grpc"
)

type ProductServiceClient struct {
	Client proto.ProductServiceClient
}

func NewProductServiceClient(url string) *ProductServiceClient {
	maxRetries := 5
	var conn *grpc.ClientConn
	var err error

	for i := 0; i < maxRetries; i++ {
		conn, err = grpc.Dial(url,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		)
		if err == nil {
			break
		}
		log.Printf("attempt %d: failed to connect to product service: %v", i+1, err)
		if i < maxRetries-1 {
			time.Sleep(time.Second * 2)
		}
	}
	if err != nil {
		log.Fatalf("failed to connect to product service after %d attempts: %v", maxRetries, err)
	}

	log.Printf("Successfully connected to product service at %s", url)
	return &ProductServiceClient{
		Client: proto.NewProductServiceClient(conn),
	}
}

func (c *ProductServiceClient) FindOne(productId int64) (*proto.FindOneResponse, error) {
	req := &proto.FindOneRequest{
		Id: productId,
	}
	t, E := c.Client.FindOne(context.Background(), req)
	log.Println(t)
	return t, E
}

func (c *ProductServiceClient) DecreaseStock(productId int64, orderId int64, quantity int64) (*proto.DecreaseStockResponse, error) {
	req := &proto.DecreaseStockRequest{
		Id:       productId,
		OrderId:  orderId,
		Quantity: quantity,
	}
	return c.Client.DecreaseStock(context.Background(), req)
}
