package repository

import (
	"context"
	"errors"

	"auth-service/internal/model"
	"auth-service/internal/utils"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Register(ctx context.Context, user *model.User) (*model.User, error) {
	var existingUser model.User
	result := r.DB.WithContext(ctx).Where("email = ?", user.Email).First(&existingUser)

	if result.Error == nil {
		return nil, errors.New("user already exists")
	}
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword
	if err := r.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UserLogin(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).Where("username = ? or email = ?", username, username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) AdminLogin(ctx context.Context, username string) (*model.Admin, error) {
	var admin model.Admin
	err := r.DB.WithContext(ctx).Where("username = ?", username).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}
	return &admin, nil
}

func (r *UserRepository) RegisterAdmin(ctx context.Context, admin *model.Admin) (*model.Admin, error) {
	hashedPassword, err := utils.HashPassword(admin.Password)
	if err != nil {
		return nil, err
	}
	admin.Password = hashedPassword

	if err := r.DB.WithContext(ctx).Create(admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}
