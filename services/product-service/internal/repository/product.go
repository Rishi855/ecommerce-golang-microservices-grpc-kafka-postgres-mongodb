package repository

import (
	"context"
	"errors"

	"product-service/internal/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	newProduct := model.Product{
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
	}
	if err := r.db.WithContext(ctx).Create(&newProduct).Error; err != nil {
		return nil, err
	}
	return &newProduct, nil
}

func (r *ProductRepository) FindOne(ctx context.Context, product *model.Product) (*model.Product, error) {
	var existing model.Product
	if err := r.db.WithContext(ctx).First(&existing, product.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &existing, nil
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) DecreaseStock(ctx context.Context, productId int64, orderId int64, quantity int64) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).First(&product, productId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// ensure idempotency: one log per orderId
	var existingLog model.StockDecreaseLog
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderId).
		First(&existingLog).Error
	if err == nil {
		return nil, errors.New("stock already decreased for this order")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if quantity <= 0 {
		return nil, errors.New("quantity must be > 0")
	}

	if product.Stock < quantity {
		return nil, errors.New("not enough stock")
	}

	product.Stock -= quantity

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	logRec := model.StockDecreaseLog{
		OrderId:      orderId,
		ProductRefer: productId,
	}
	if err := tx.Create(&logRec).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &product, nil
}
