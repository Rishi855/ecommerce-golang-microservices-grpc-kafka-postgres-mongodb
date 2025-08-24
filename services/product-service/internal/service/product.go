package service

import (
	"context"
	"errors"

	"product-service/internal/model"
	"product-service/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	if product == nil {
		return nil, errors.New("product is required")
	}
	if product.Name == "" || product.Price <= 0 || product.Stock < 0 {
		return nil, errors.New("invalid product data")
	}
	return s.repo.CreateProduct(ctx, product)
}

func (s *ProductService) FindOne(ctx context.Context, product *model.Product) (*model.Product, error) {
	if product == nil || product.Id <= 0 {
		return nil, errors.New("invalid product id")
	}
	return s.repo.FindOne(ctx, product)
}

func (s *ProductService) FindAll(ctx context.Context) ([]model.Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *ProductService) DecreaseStock(ctx context.Context, productId int64, orderId int64, quantity int64) (*model.Product, error) {
	if productId <= 0 || orderId <= 0 || quantity <= 0 {
		return nil, errors.New("invalid input for decrease stock")
	}
	return s.repo.DecreaseStock(ctx, productId, orderId, quantity)
}
