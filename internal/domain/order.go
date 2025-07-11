package domain

import "time"

type orderStatus string

const Pending orderStatus = "PENDING"
const Paid orderStatus = "PAID"
const Cancelled orderStatus = "CANCELLED"

type Order struct {
	ID         int64        `json:"id"`
	UserID     int64        `json:"user_id"`
	TotalPrice float64      `json:"total_price"`
	Discount   float64      `json:"discount"`
	FinalPrice float64      `json:"final_price"`
	CouponCode string       `json:"coupon_code"`
	Status     orderStatus  `json:"status"`
	Items      []*OrderItem `json:"items"`
	CreatedAt  time.Time    `json:"created_at"`
}

type OrderItem struct {
	ID          int64   `json:"id"`
	OrderID     int64   `json:"order_id"`
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}
