package repository

type OrderRepository struct{
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository{
	return &OrderRepository{db:db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *model.Order)