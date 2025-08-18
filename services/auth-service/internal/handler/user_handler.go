package handler

import (
	"context"
	"net/http"

	"auth-service/internal/model"
	"auth-service/internal/proto"
	"auth-service/internal/service"
	"auth-service/internal/utils"
)

type AuthHandler struct {
	proto.UnimplementedAuthServiceServer
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := h.service.Register(ctx, user)
	if err != nil {
		return &proto.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	return &proto.RegisterResponse{
		Status: http.StatusCreated,
		Id:     res.Id,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	res, err := h.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		return &proto.LoginResponse{
			Status: http.StatusUnauthorized,
			Error:  "Invalid email or password",
		}, nil
	}

	token, err := utils.GenerateToken(res)
	if err != nil {
		return &proto.LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}
	return &proto.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (h *AuthHandler) AdminRegister(ctx context.Context, req *proto.AdminRegisterRequest) (*proto.RegisterResponse, error) {
	admin := &model.Admin{
		Username: req.Username,
		Password: req.Password,
	}

	res, err := h.service.AdminRegister(ctx, admin)
	if err != nil {
		return &proto.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	return &proto.RegisterResponse{
		Status: http.StatusCreated,
		Id:     res.Id,
	}, nil
}

func (h *AuthHandler) AdminLogin(ctx context.Context, req *proto.AdminLoginRequest) (*proto.LoginResponse, error) {
	res, err := h.service.AdminLogin(ctx, req.Username)
	if err != nil {
		return &proto.LoginResponse{
			Status: http.StatusUnauthorized,
			Error:  "User not found",
		}, nil
	}

	isPasswordMatch := utils.CheckHashPassword(req.Password, res.Password)
	if !isPasswordMatch {
		return &proto.LoginResponse{
			Status: http.StatusUnauthorized,
			Error:  "Invalid username or password",
		}, nil
	}

	token, err := utils.GenerateAdminToken(res)
	if err != nil {
		return &proto.LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}
	return &proto.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (h *AuthHandler) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	claims, err := utils.ValidateToken(req.Token)
	if err != nil {
		return &proto.ValidateResponse{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		}, nil
	}

	return &proto.ValidateResponse{
		Status: http.StatusOK,
		UserID: claims.Id,
	}, nil
}
