package handler

import (
	"context"
	"net/http"

	"order-service/internal/model"
	"order-service/internal/proto"

	"gorm.io/gorm"
)

type OrderHandler struct {
	db            *gorm.DB
	productClient proto.ProductServiceClient
	proto.UnimplementedOrderServiceServer
}

func NewOrderHandler(db *gorm.DB, productClient proto.ProductServiceClient) *OrderHandler {
	return &OrderHandler{
		db:            db,
		productClient: productClient,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	product, err := h.productClient.FindOne(ctx, &proto.FindOneRequest{Id: req.ProductId})
	if err != nil {
		return &proto.CreateOrderResponse{
			Status: http.StatusBadGateway,
			Error:  err.Error(),
		}, nil
	}

	if product.Status >= http.StatusNotFound {
		return &proto.CreateOrderResponse{
			Status: product.Status,
			Error:  product.Error,
		}, nil
	}

	if product.Data.Stock <= 0 {
		return &proto.CreateOrderResponse{
			Status: http.StatusBadRequest,
			Error:  "Stock is low",
		}, nil
	}

	order := model.Order{
		Price:     product.Data.Price,
		ProductId: product.Data.Id,
		UserId:    req.UserId,
	}

	if err := h.db.Create(&order).Error; err != nil {
		return &proto.CreateOrderResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	res, err := h.productClient.DecreaseStock(ctx, &proto.DecreaseStockRequest{
		Id:      order.ProductId,
		OrderId: order.Id,
		Quantity: req.Quantity,
	})

	if err != nil {
		return &proto.CreateOrderResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	if res.Status == "409" { // conflict
		h.db.Delete(&model.Order{}, order.Id)
		return &proto.CreateOrderResponse{
			Status: http.StatusConflict,
			Error:  res.Error,
		}, nil
	}

	return &proto.CreateOrderResponse{
		Status: http.StatusCreated,
		Id:     order.Id,
	}, nil
}
