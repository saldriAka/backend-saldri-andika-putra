package domain

import (
	"context"
	"saldri/backend-saldri-andika-putra/dto"
	"time"
)

type Transaction struct {
	ID          int    `gorm:"primaryKey;autoIncrement"`
	CustomerID  string `gorm:"not null"`
	Customer    Users  `gorm:"foreignKey:CustomerID;references:ID"`
	CreatedAt   time.Time
	ShippingFee float64
	Discount    float64
	FinalAmount float64
	TotalAmount float64
	Items       []TransactionItem `gorm:"foreignKey:TransactionID"`
}

type TransactionItem struct {
	ID            int     `gorm:"primaryKey;autoIncrement"`
	TransactionID int     `gorm:"not null"` // akan diisi GORM otomatis karena relasi
	ProductID     int     `gorm:"not null"`
	Quantity      int     `gorm:"not null"`
	UnitPrice     float64 `gorm:"not null"`
	TotalPrice    float64 `gorm:"not null"`

	Product Product `gorm:"foreignKey:ProductID;references:ID"`
}

type TransactionRepository interface {
	Create(ctx context.Context, tx Transaction) error
	FindByUserID(ctx context.Context, userID string) ([]Transaction, error)
	FindCustomersByMerchantID(ctx context.Context, merchantID string) ([]Users, error)
}

type TransactionService interface {
	Create(ctx context.Context, req dto.CreateTransactionRequest, userID string) error
	GetUserTransactions(ctx context.Context, userID string) ([]Transaction, error)
	GetCustomersByMerchantID(ctx context.Context, merchantID string) ([]Users, error)
}
