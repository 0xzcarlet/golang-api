package product

import "time"

type CreateProductReq struct {
	Name  string `json:"name" validate:"required,min=2,max=120"`
	Price int    `json:"price" validate:"required,min=0,max=2000000000"`
}

type Product struct {
	ID        uint64    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Price     int       `db:"price" json:"price"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type ErrResp struct {
	Error string `json:"error"`
}
type UpdateProductReq struct {
	Name  *string `json:"name" validate:"omitempty,min=2,max=120"`
	Price *int    `json:"price" validate:"omitempty,min=0,max=2000000000"`
}
