package auth

import (
	"api-gateaway/internal/auth/proto"

	"google.golang.org/grpc"
)

type AuthClient struct {
	AuthClient proto.AuthServiceClient
}

func NewAuthServiceClient(conn *grpc.ClientConn) *AuthClient {
	return &AuthClient{AuthClient: proto.NewAuthServiceClient(conn)}
}
