package handler

import (
	"context"
	"net/http"

	"order-service/internal/model"
	"order-service/internal/proto"
	"order-service/internal/client"

	"gorm.io/gorm"
)

type OrderHandler struct {
	db            *gorm.DB
	ProductClient *client.ProductServiceClient
	proto.UnimplementedOrderServiceServer
}

func NewOrderHandler(db *gorm.DB, productClient *client.ProductServiceClient) *OrderHandler {
	return &OrderHandler{
		db:            db,
		ProductClient: productClient,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {

	productResp, err := h.ProductClient.FindOne(req.ProductId)
	if err != nil {
		return &proto.CreateOrderResponse{
			Status: int64(http.StatusBadGateway),
			Error:  err.Error(),
		}, nil
	}

	if productResp.Status >= int64(http.StatusNotFound) {
		return &proto.CreateOrderResponse{
			Status: productResp.Status,
			Error:  productResp.Error,
		}, nil
	}

	if productResp.Data == nil || productResp.Data.Stock <= 0 {
		return &proto.CreateOrderResponse{
			Status: int64(http.StatusBadRequest),
			Error:  "Stock is low",
		}, nil
	}

	order := model.Order{
		Price:     productResp.Data.Price,
		ProductId: productResp.Data.Id,
		UserId:    req.UserId,
	}

	if err := h.db.Create(&order).Error; err != nil {
		return &proto.CreateOrderResponse{
			Status: int64(http.StatusInternalServerError),
			Error:  err.Error(),
		}, nil
	}

	decResp, err := h.ProductClient.DecreaseStock(order.ProductId,order.Id,req.Quantity)
	if err != nil {		_ = h.db.Delete(&model.Order{}, order.Id).Error
		return &proto.CreateOrderResponse{
			Status: int64(http.StatusBadRequest),
			Error:  err.Error(),
		}, nil
	}

	if decResp.Status == "409" { 
		_ = h.db.Delete(&model.Order{}, order.Id).Error
		return &proto.CreateOrderResponse{
			Status: int64(http.StatusConflict),
			Error:  decResp.Error,
		}, nil
	}

	return &proto.CreateOrderResponse{
		Status: int64(http.StatusCreated),
		Id:     order.Id,
	}, nil
}
