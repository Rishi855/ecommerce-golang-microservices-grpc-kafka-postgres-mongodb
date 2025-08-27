package client

import (
	"context"
	"log"

	"order-service/internal/proto"

	"google.golang.org/grpc"
)

type ProductServiceClient struct {
	Client proto.ProductServiceClient
}

func NewProductServiceClient(url string) *ProductServiceClient {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to product service: %v", err)
	}
	return &ProductServiceClient{
		Client: proto.NewProductServiceClient(conn),
	}
}

func (c *ProductServiceClient) FindOne(productId int64) (*proto.FindOneResponse, error){
	req := &proto.FindOneRequest{
		Id : productId,
	}
	t,E := c.Client.FindOne(context.Background(),req)
	log.Println(t)
	return t,E
}

func (c *ProductServiceClient) DecreaseStock(productId int64, orderId int64, quantity int64) (*proto.DecreaseStockResponse, error){
	req :=  &proto.DecreaseStockRequest{
		Id : productId,
		OrderId : orderId,
		Quantity: quantity,
	}
	return c.Client.DecreaseStock(context.Background(), req)
}