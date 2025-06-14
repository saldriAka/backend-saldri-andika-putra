package dto

type CreateProductRequest struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	Name        string  `json:"name" validate:"required,min=3"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Description string  `json:"description" validate:"required,min=5"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	MerchantID  string  `json:"merchant_id" validate:"required,uuid4"`
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
	Name        string  `json:"name" validate:"omitempty,min=3"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
	Description string  `json:"description" validate:"omitempty,min=5"`
	Stock       int     `json:"stock" validate:"omitempty,gte=0"`
}
