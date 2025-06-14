package domain

import (
	"context"
	"saldri/backend-saldri-andika-putra/dto"
	"time"
)

type Users struct {
	ID        string `gorm:"primaryKey;type:char(36)"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"not null"`
	CreatedAt time.Time
}

type UsersService interface {
	Register(ctx context.Context, req dto.RegisterUserRequest) error
	Login(context.Context, dto.AuthRequest) (dto.AuthResponse, error)
	GetProfile(ctx context.Context, id string) (dto.UserData, error)
	List(ctx context.Context, page, limit int) ([]dto.UserData, int64, error)
}

type UsersRepository interface {
	Save(ctx context.Context, user *Users) error
	FindByEmail(ctx context.Context, email string) (Users, error)
	FindById(ctx context.Context, id string) (Users, error)
	FindAll(ctx context.Context, limit, offset int) ([]Users, int64, error)
}
