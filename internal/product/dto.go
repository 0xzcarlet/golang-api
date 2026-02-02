package product

// Request DTOs
type CreateProductReq struct {
	Name  string `json:"name" validate:"required,min=2,max=120"`
	Price int    `json:"price" validate:"required,min=0,max=2000000000"`
}

type UpdateProductReq struct {
	Name  *string `json:"name" validate:"omitempty,min=2,max=120"`
	Price *int    `json:"price" validate:"omitempty,min=0,max=2000000000"`
}

// Response DTOs
type ProductResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	CreatedAt string `json:"created_at"`
}
