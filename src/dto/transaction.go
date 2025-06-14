package dto

import "time"

type CreateTransactionRequest struct {
	ProductID int `json:"product_id" validate:"required,gt=0"`
	Quantity  int `json:"quantity" validate:"required,gt=0"`
}

type Transaction struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	UserID    string
	User      UserData `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	Items     []TransactionItem `gorm:"foreignKey:TransactionID"` // Tambahkan ini
}

type TransactionItem struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	TransactionID int
	ProductID     int
	Quantity      int
	Price         float64

	Product Product `gorm:"foreignKey:ProductID"`
}
