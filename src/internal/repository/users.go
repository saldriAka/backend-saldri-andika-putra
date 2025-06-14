package repository

import (
	"context"
	"saldri/backend-saldri-andika-putra/domain"

	"gorm.io/gorm"
)

type usersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) domain.UsersRepository {
	return &usersRepository{db: db}
}

func (r *usersRepository) Save(ctx context.Context, user *domain.Users) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *usersRepository) FindByEmail(ctx context.Context, email string) (domain.Users, error) {
	var user domain.Users
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *usersRepository) FindById(ctx context.Context, id string) (domain.Users, error) {
	var user domain.Users
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

func (r *usersRepository) FindAll(ctx context.Context, limit, offset int) ([]domain.Users, int64, error) {
	var users []domain.Users
	var total int64

	db := r.db.WithContext(ctx).Model(&domain.Users{})
	db.Count(&total)

	err := db.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}
