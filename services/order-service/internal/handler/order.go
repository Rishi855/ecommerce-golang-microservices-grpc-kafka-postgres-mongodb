package handler

import (
	"context"
	"encoding/json"
	"time"

	// "log"
	"net/http"

	"order-service/internal/client"
	"order-service/internal/events"
	"order-service/internal/model"
	"order-service/internal/proto"

	"gorm.io/gorm"
)

type OrderHandler struct {
	db                   *gorm.DB
	ProductClient        *client.ProductServiceClient
	NotificationProducer *events.Producer
	LogProducer          *events.Producer
	proto.UnimplementedOrderServiceServer
}

func NewOrderHandler(db *gorm.DB, productClient *client.ProductServiceClient, np *events.Producer, lp *events.Producer) *OrderHandler {
	return &OrderHandler{
		db:                   db,
		ProductClient:        productClient,
		NotificationProducer: np,
		LogProducer:          lp,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {

	productResp, err := h.ProductClient.FindOne(req.ProductId)
	if err != nil {
		h.publishLog("ERROR", "Product service FindOne failed: "+err.Error())
		return &proto.CreateOrderResponse{
			Status: int64(http.StatusBadGateway),
			Error:  err.Error(),
		}, nil
	}

	if productResp.Status >= int64(http.StatusNotFound) {
		h.publishLog("WARN", "Product not found for id: "+string(rune(req.ProductId)))
		return &proto.CreateOrderResponse{
			Status: productResp.Status,
			Error:  productResp.Error,
		}, nil
	}
	if productResp.Data == nil || productResp.Data.Stock < req.Quantity {
		h.publishLog("WARN", "Not enough stock for product")
		return &proto.CreateOrderResponse{
			Status: int64(http.StatusBadRequest),
			Error:  "Not enough stock",
		}, nil
	}

	order := model.Order{
		Price:     productResp.Data.Price,
		ProductId: productResp.Data.Id,
		UserId:    req.UserId,
	}

	if err := h.db.Create(&order).Error; err != nil {
		h.publishLog("ERROR", "DB insert failed: "+err.Error())
		return &proto.CreateOrderResponse{
			Status: int64(http.StatusInternalServerError),
			Error:  err.Error(),
		}, nil
	}

	decResp, err := h.ProductClient.DecreaseStock(order.ProductId, order.Id, req.Quantity)
	if err != nil {
		h.publishLog("ERROR", "Decrease stock failed: "+err.Error())
		_ = h.db.Delete(&model.Order{}, order.Id).Error
		return &proto.CreateOrderResponse{
			Status: int64(http.StatusBadRequest),
			Error:  err.Error(),
		}, nil
	}

	if decResp.Status == http.StatusConflict {
		_ = h.db.Delete(&model.Order{}, order.Id).Error
		h.publishLog("WARN", "Decrease stock conflict for order")
		return &proto.CreateOrderResponse{
			Status: int64(http.StatusConflict),
			Error:  decResp.Error,
		}, nil
	}

	event := map[string]interface{}{
		"event":    "notification.trigger",
		"user_id":  order.UserId,
		"channel":  "email",
		"template": "ORDER_CONFIRMED",
		"data": map[string]interface{}{
			"order_id":     order.Id,
			"amount":       order.Price,
			"product_name": productResp.Data.Name,
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	data, _ := json.Marshal(event)
	_ = h.NotificationProducer.Publish("orderKey",string(data))

	return &proto.CreateOrderResponse{
		Status: int64(http.StatusCreated),
		Id:     order.Id,
	}, nil
}

func (h *OrderHandler) publishLog(level, message string) {
	logEvent := map[string]interface{}{
		"level":     level,
		"service":   "order-service",
		"message":   message,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	data, _ := json.Marshal(logEvent)
	_ = h.LogProducer.Publish("logKey", string(data))
}
