package service

import (
	"context"
	"errors"
	"saldri/backend-saldri-andika-putra/domain"
	"saldri/backend-saldri-andika-putra/dto"
)

type productService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) domain.ProductService {
	return &productService{repo}
}

func (s *productService) GetAll(ctx context.Context) ([]domain.Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *productService) GetByID(ctx context.Context, id int) (domain.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *productService) Create(ctx context.Context, req dto.CreateProductRequest) error {
	product := domain.Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Stock:       req.Stock,
		MerchantID:  req.MerchantID,
	}
	return s.repo.Create(ctx, product)
}

func (s *productService) UpdateStock(ctx context.Context, id int, stock int) error {
	return s.repo.UpdateStock(ctx, id, stock)
}

func (s *productService) Update(ctx context.Context, id int, req dto.UpdateProductRequest, merchantID string) error {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if product.MerchantID != merchantID {
		return errors.New("tidak diizinkan mengubah produk merchant lain")
	}

	product.Name = req.Name
	product.Price = req.Price
	product.Description = req.Description
	product.Stock = req.Stock

	return s.repo.Update(ctx, product)
}

func (s *productService) Delete(ctx context.Context, id int, merchantID string) error {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if product.MerchantID != merchantID {
		return errors.New("tidak diizinkan menghapus produk merchant lain")
	}

	return s.repo.Delete(ctx, id)
}
