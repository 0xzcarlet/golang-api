package place

import "time"

// Place DTOs

type CreatePlaceReq struct {
	Name        *string    `json:"name" validate:"omitempty,max=255"`
	Link        *string    `json:"link" validate:"omitempty,max=255"`
	LinkType    *int       `json:"link_type" validate:"omitempty,min=0"`
	Description *string    `json:"description" validate:"omitempty"`
	GoAt        *time.Time `json:"go_at" validate:"omitempty"`
	GoAtTime    *time.Time `json:"go_at_time" validate:"omitempty"`
	Status      *int       `json:"status" validate:"omitempty,min=0"`
}

type UpdatePlaceReq struct {
	Name        *string    `json:"name" validate:"omitempty,max=255"`
	Link        *string    `json:"link" validate:"omitempty,max=255"`
	LinkType    *int       `json:"link_type" validate:"omitempty,min=0"`
	Description *string    `json:"description" validate:"omitempty"`
	GoAt        *time.Time `json:"go_at" validate:"omitempty"`
	GoAtTime    *time.Time `json:"go_at_time" validate:"omitempty"`
	Status      *int       `json:"status" validate:"omitempty,min=0"`
}

type PlaceResponse struct {
	ID          uint64     `json:"id"`
	UserID      uint64     `json:"user_id"`
	Name        *string    `json:"name"`
	Link        *string    `json:"link"`
	LinkType    *int       `json:"link_type"`
	Description *string    `json:"description"`
	GoAt        *time.Time `json:"go_at"`
	GoAtTime    *time.Time `json:"go_at_time"`
	Status      *int       `json:"status"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

// PlaceCategory DTOs

type CreatePlaceCategoryReq struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}

type UpdatePlaceCategoryReq struct {
	Name *string `json:"name" validate:"omitempty,min=1,max=50"`
}

type PlaceCategoryResponse struct {
	ID     uint   `json:"id"`
	UserID uint64 `json:"user_id"`
	Name   string `json:"name"`
}
