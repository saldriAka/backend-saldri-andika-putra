package domain

import (
	"context"
	"saldri/backend-saldri-andika-putra/dto"
)

type Product struct {
	ID          int `gorm:"primaryKey"`
	Name        string
	Price       float64
	Description string
	Stock       int
	MerchantID  string `gorm:"not null"`
}

type ProductRepository interface {
	FindAll(ctx context.Context) ([]Product, error)
	FindByID(ctx context.Context, id int) (Product, error)
	Create(ctx context.Context, product Product) error
	UpdateStock(ctx context.Context, id int, stock int) error
	Update(ctx context.Context, product Product) error
	Delete(ctx context.Context, id int) error
}

type ProductService interface {
	GetAll(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id int) (Product, error)
	Create(ctx context.Context, req dto.CreateProductRequest) error
	Update(ctx context.Context, productID int, updated dto.UpdateProductRequest, merchantID string) error
	Delete(ctx context.Context, productID int, merchantID string) error
}
