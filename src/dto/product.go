package dto

type CreateProductRequest struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Stock       int     `json:"stock"`
	MerchantID  string  `json:"-"`
}

type Product struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	Name        string
	Price       float64
	Description string
	Stock       int
	MerchantID  string
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Stock       int     `json:"stock"`
}
