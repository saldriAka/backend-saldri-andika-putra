package repository

import (
	"context"
	"saldri/backend-saldri-andika-putra/domain"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) FindByUserID(ctx context.Context, userID string) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	err := r.db.WithContext(ctx).
		Preload("Items.Product").
		Preload("Customer"). // ⬅️ Tambahkan ini
		Where("customer_id = ?", userID).
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepository) FindCustomersByMerchantID(ctx context.Context, merchantID string) ([]domain.Users, error) {
	var users []domain.Users

	err := r.db.WithContext(ctx).
		Model(&domain.TransactionItem{}).
		Select("DISTINCT users.*").
		Joins("JOIN transactions ON transaction_items.transaction_id = transactions.id").
		Joins("JOIN products ON transaction_items.product_id = products.id").
		Joins("JOIN users ON transactions.customer_id = users.id").
		Where("products.merchant_id = ?", merchantID).
		Scan(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *transactionRepository) Create(ctx context.Context, tx domain.Transaction) error {
	return r.db.WithContext(ctx).Transaction(func(txDb *gorm.DB) error {
		if err := txDb.Create(&tx).Error; err != nil {
			return err
		}
		return nil
	})
}
