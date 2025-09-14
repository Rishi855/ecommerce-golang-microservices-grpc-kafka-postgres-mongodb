package model

type Order struct {
	Id        int64 `json:"id" gorm:"primaryKey"`
	ProductId int64 `json:"product_id"`
	Price     int64 `json:"price"`
	UserId    int64 `json:"user_id"`
}