package auth

import (
	"api-gateway/internal/auth/proto"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Client *AuthClient
}

func NewAuthHandler(client *AuthClient) *AuthHandler {
	return &AuthHandler{Client: client}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := h.Client.ServiceClient.Register(context.Background(), &proto.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(int(res.Status), res)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req struct {
		Username    string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := h.Client.ServiceClient.Login(context.Background(), &proto.LoginRequest{
		Username:    req.Username,
		Password: req.Password,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.SetCookie("Authorization", res.Token, 3600*24*30, "", "", false, true)
	ctx.JSON(http.StatusOK, res)
}

func (h *AuthHandler) AdminRegister(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := h.Client.ServiceClient.AdminRegister(context.Background(), &proto.AdminRegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(int(res.Status), res)
}

func (h *AuthHandler) AdminLogin(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := h.Client.ServiceClient.AdminLogin(context.Background(), &proto.AdminLoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(int(res.Status), res)
}
