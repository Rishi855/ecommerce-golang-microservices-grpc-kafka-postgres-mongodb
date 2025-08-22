package repository

import (
	"context"
	"errors"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, err) {
	newProduct := model.Product{
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
	}
	err = r.db.WithContext(ctx).Create(&newProduct).Error
	if err != nil {
		return nil, err
	}
	return &newProduct, nil
}

func (r *ProductRepository) FindOne(ctx context.Context, product *model.Product)(*model.Product, error){
	var existingProduct model.Product

	err := r.db.WithContext(ctx).First(&existingProduct,product.Id).Error
	if err!=nil{
		if errors.Is(err,gorm.ErrorRecordNotFound){
			return nil,errors.New("Product not found")
		}
		return nil,err
	}
	return &existingProduct,nil
}

func (r *ProductRepository) FindAll(ctx context.Context)([]model.Product, err){
	var products []model.Product

	err := r.db.WithContext(ctx).First(&existingProduct, product.Id).Error
	if err!=nil{
		return nil,err
	}
	return products nil
}

func (r *ProductRepository) DecreaseStock(ctx context.Context, productId int64, orderId int64,quantity int64) (*model.Produt, err){
	
}