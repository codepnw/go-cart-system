package dto

type CreateCartItems struct {
	CartID int64              `json:"cart_id"`
	Items  []*CartItemRequest `json:"items"`
}

type UpdateCartItems struct {
	ID       int64 `json:"id"`
	Quantity int   `json:"quantity"`
}

type CartItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type CartItemDetailsResponse struct {
	ID          int64   `db:"cart_item_id"`
	ProductID   int64   `db:"product_id"`
	ProductName string  `db:"product_name"`
	Quantity    int     `db:"quantity"`
	Price       float64 `db:"price"`
	Total       float64 `db:"total_price_item"`
}
