package auth

import (
	"api-gateway/internal/auth/proto"

	"google.golang.org/grpc"
)

type AuthClient struct {
	ServiceClient proto.AuthServiceClient
}

func NewAuthServiceClient(conn *grpc.ClientConn) *AuthClient {
	return &AuthClient{ServiceClient: proto.NewAuthServiceClient(conn)}
}
