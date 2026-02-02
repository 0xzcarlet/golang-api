package place

import (
	"database/sql"
	"time"
)

// Place represents the place domain model
type Place struct {
	ID          uint64         `db:"id" json:"id"`
	UserID      uint64         `db:"user_id" json:"user_id"`
	Name        sql.NullString `db:"name" json:"name"`
	Link        sql.NullString `db:"link" json:"link"`
	LinkType    sql.NullInt32  `db:"link_type" json:"link_type"`
	Description sql.NullString `db:"description" json:"description"`
	GoAt        sql.NullTime   `db:"go_at" json:"go_at"`
	GoAtTime    sql.NullTime   `db:"go_at_time" json:"go_at_time"`
	Status      sql.NullInt32  `db:"status" json:"status"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

// PlaceCategory represents the place_category domain model
type PlaceCategory struct {
	ID     uint   `db:"id" json:"id"`
	UserID uint64 `db:"user_id" json:"user_id"`
	Name   string `db:"name" json:"name"`
}
