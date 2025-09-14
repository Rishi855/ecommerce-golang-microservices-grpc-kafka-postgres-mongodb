package handler

import (
	"context"
	"log"
	"net/http"

	"product-service/internal/model"
	"product-service/internal/proto"
	"product-service/internal/service"
)

type ProductHandler struct {
	proto.UnimplementedProductServiceServer
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	input := &model.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	created, err := h.svc.CreateProduct(ctx, input)
	if err != nil {
		return &proto.CreateProductResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
			Id:     0,
		}, nil
	}

	return &proto.CreateProductResponse{
		Status: http.StatusCreated,
		Error:  "",
		Id:     created.Id,
	}, nil
}

func (h *ProductHandler) FindOne(ctx context.Context, req *proto.FindOneRequest) (*proto.FindOneResponse, error) {
	found, err := h.svc.FindOne(ctx, &model.Product{Id: req.Id})
	if err != nil {
		return &proto.FindOneResponse{
			Status: http.StatusNotFound,
			Error:  err.Error(),
			Data:   nil,
		}, nil
	}
	t := proto.FindOneData{
			Id:    found.Id,
			Name:  found.Name,
			Price: found.Price,
			Stock: found.Stock,
		}
	log.Println("here product handler: ",found)
	return &proto.FindOneResponse{
		Status: http.StatusOK,
		Error:  "",
		Data: &t,
	}, nil
}

func (h *ProductHandler) FindAll(ctx context.Context, req *proto.FindAllRequest) (*proto.FindAllResponse, error) {
	items, err := h.svc.FindAll(ctx)
	if err != nil {
		return &proto.FindAllResponse{
			Status:   http.StatusBadGateway,
			Error:    err.Error(),
			Products: nil,
		}, nil
	}

	var out []*proto.FindOneData
	for _, p := range items {
		p := p
		out = append(out, &proto.FindOneData{
			Id:    p.Id,
			Name:  p.Name,
			Price: p.Price,
			Stock: p.Stock,
		})
	}

	return &proto.FindAllResponse{
		Status:   http.StatusOK,
		Error:    "",
		Products: out,
	}, nil
}

func (h *ProductHandler) DecreaseStock(ctx context.Context, req *proto.DecreaseStockRequest) (*proto.DecreaseStockResponse, error) {
	_, err := h.svc.DecreaseStock(ctx, req.Id, req.OrderID, req.Quantity)
	if err != nil {
		return &proto.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  err.Error(),
		}, nil
	}

	return &proto.DecreaseStockResponse{
		Status: http.StatusOK,
		Error:  "",
	}, nil
}
