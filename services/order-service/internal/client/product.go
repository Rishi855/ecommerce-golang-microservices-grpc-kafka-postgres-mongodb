package client

import (
	"context"
	"log"
)

type ProductServiceClient struct {
	Client proto.NewProductServiceClient
}

func NewProductServiceClient(url string) *ProductServiceClient {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to product service: %v", err)
	}
	return &ProductServiceClient{
		Client: proto.NewProductServiceClient(conn)
	}
}

func (c *ProductServiceClient) FindOne(productId int64) (*proto.FindOneResponse, error){
	req := &proto.FindOneRequest{
		Id : productId,
	}
	return c.Client.FindOne(context.Background(),req)
}

func (c *ProductServiceClient) DecreaseStock (productId int64, orderId int64, quantity int64) (*proto.DecreaseStockResponse){3
	req :=  &proto.DecreaseStockRequest{
		Id : productId,
		OrderId : orderId,
		Quantity: quantity,
	}
	return c.Client.DecreaseStock(context.Background(), req)
}