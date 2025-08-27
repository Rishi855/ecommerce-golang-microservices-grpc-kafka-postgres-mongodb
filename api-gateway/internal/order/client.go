package order

import (
	"api-gateway/internal/order/proto"

	"google.golang.org/grpc"
)

type OrderClient struct {
	OrderService proto.OrderServiceClient
}

func NewOrderServiceClient(conn *grpc.ClientConn) *OrderClient{
	return &OrderClient{
		OrderService: proto.NewOrderServiceClient(conn),
	}
}