package auth

import (
	"api-gateway/internal/auth/proto"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	Client *AuthClient
}

func NewAuthMiddleware(client *AuthClient) *AuthMiddleware {
	return &AuthMiddleware{Client: client}
}

func (s *AuthMiddleware) ValidateToken(ctx *gin.Context) {
	authorization, err := ctx.Cookie("Authorization")
	if err != nil || authorization == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization cookie"})
		return
	}

	token := strings.TrimSpace(authorization)

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is empty"})
		return
	}

	res, err := s.Client.ServiceClient.Validate(context.Background(), &proto.ValidateRequest{
		Token: token,
	})

	if err != nil || res.Status != http.StatusOK {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	ctx.Set("userId", res.UserID)
	ctx.Next()
}
