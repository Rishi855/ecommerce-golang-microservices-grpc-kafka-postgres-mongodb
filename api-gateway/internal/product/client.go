package product

import (
	"api-gateway/internal/product/proto"

	"google.golang.org/grpc"
)

type ProductClient struct {
	ProductService proto.ProductServiceClient
}

func NewProductServiceClient(conn *grpc.ClientConn) *ProductClient{
	return &ProductClient{
		ProductService: proto.NewProductServiceClient(conn),
	}
}