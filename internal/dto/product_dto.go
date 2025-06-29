package dto

type CreateProduct struct {
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"required"`
}

type UpdateProduct struct {
	ID    int64    `json:"id"`
	Name  *string  `json:"name"`
	Price *float64 `json:"price"`
}
