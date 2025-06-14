package service

import (
	"context"
	"errors"
	"saldri/backend-saldri-andika-putra/domain"
	"saldri/backend-saldri-andika-putra/dto"
	"time"
)

type transactionService struct {
	txRepo   domain.TransactionRepository
	prodRepo domain.ProductRepository
}

func NewTransactionService(txRepo domain.TransactionRepository, prodRepo domain.ProductRepository) domain.TransactionService {
	return &transactionService{txRepo, prodRepo}
}

func (s *transactionService) Create(ctx context.Context, req dto.CreateTransactionRequest, userID string) error {
	// Ambil data produk
	product, err := s.prodRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		return err
	}

	if product.Stock < req.Quantity {
		return errors.New("stok produk tidak mencukupi")
	}

	// Hitung total awal
	unitPrice := product.Price
	total := unitPrice * float64(req.Quantity)

	discount := 0.0
	shippingFee := 10000.0 // default ongkir

	if total > 15000 {
		shippingFee = 0
	}
	if total > 50000 {
		discount = total * 0.10
	}

	finalTotal := total - discount + shippingFee

	// Buat transaksi
	transaction := domain.Transaction{
		CustomerID:  userID,
		CreatedAt:   time.Now(),
		ShippingFee: shippingFee,
		Discount:    discount,
		FinalAmount: finalTotal,
		TotalAmount: total,
		Items: []domain.TransactionItem{
			{
				ProductID:  product.ID,
				Quantity:   req.Quantity,
				UnitPrice:  unitPrice,
				TotalPrice: total,
			},
		},
	}

	// Simpan transaksi ke database
	if err := s.txRepo.Create(ctx, transaction); err != nil {
		return err
	}

	// Update stok produk
	newStock := product.Stock - req.Quantity
	if err := s.prodRepo.UpdateStock(ctx, product.ID, newStock); err != nil {
		return err
	}

	return nil
}

func (s *transactionService) GetUserTransactions(ctx context.Context, userID string) ([]domain.Transaction, error) {
	return s.txRepo.FindByUserID(ctx, userID)
}

func (s *transactionService) GetCustomersByMerchantID(ctx context.Context, merchantID string) ([]domain.Users, error) {
	return s.txRepo.FindCustomersByMerchantID(ctx, merchantID)
}
