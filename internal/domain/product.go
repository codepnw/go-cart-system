package domain

import "time"

type Product struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Price     float64    `json:"price"`
	Stock     int        `json:"stock"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
