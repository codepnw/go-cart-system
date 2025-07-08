package dto

type CreateCartItems struct {
	CartID int64              `json:"cart_id" validate:"required"`
	Items  []*CartItemRequest `json:"items" validate:"required,min=1"`
}

type UpdateCartItem struct {
	ProductID int64 `json:"product_id" validate:"required"`
	Quantity  int   `json:"quantity" validate:"required,min=0"`
}

type CartItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type CartItem struct {
	ID          int64   `db:"cart_item_id"`
	ProductID   int64   `db:"product_id"`
	ProductName string  `db:"product_name"`
	Quantity    int     `db:"quantity"`
	Price       float64 `db:"price"`
	Total       float64 `db:"total_price_item"`
}

type CartResponse struct {
	Items      []*CartItem `json:"items"`
	TotalItems int         `json:"total_items"`
	TotalPrice float64     `json:"total_price"`
}
