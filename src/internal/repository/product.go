package repository

import (
	"context"
	"saldri/backend-saldri-andika-putra/domain"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) FindAll(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) FindByID(ctx context.Context, id int) (domain.Product, error) {
	var product domain.Product
	if err := r.db.WithContext(ctx).First(&product, "id = ?", id).Error; err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (r *productRepository) Create(ctx context.Context, product domain.Product) error {
	return r.db.WithContext(ctx).Create(&product).Error
}

func (r *productRepository) UpdateStock(ctx context.Context, id int, stock int) error {
	return r.db.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ?", id).
		Update("stock", stock).Error
}

func (r *productRepository) Update(ctx context.Context, product domain.Product) error {
	return r.db.WithContext(ctx).Save(&product).Error
}

// ✅ Hapus produk
func (r *productRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, id).Error
}
