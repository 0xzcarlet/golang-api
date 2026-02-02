package product

import "time"

// Product represents the product domain model
type Product struct {
	ID        uint64    `db:"id" json:"id"`
	UserID    uint64    `db:"user_id" json:"user_id"`
	Name      string    `db:"name" json:"name"`
	Price     int       `db:"price" json:"price"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
