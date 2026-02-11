package place

import (
	"context"
	"go-saas-api/pkg/customtime"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Place Repository Methods

func (r *Repository) CreatePlace(ctx context.Context, userID uint64, req CreatePlaceReq) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO place (user_id, name, link, link_type, description, go_at, go_at_time, status) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, req.Name, req.Link, req.LinkType, req.Description, req.GoAt, req.GoAtTime, req.Status,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) ListPlaces(ctx context.Context, userID uint64, limit int) ([]Place, error) {
	var items []Place
	err := r.db.SelectContext(ctx, &items,
		`SELECT id, user_id, name, link, link_type, description, go_at, go_at_time, status, created_at, updated_at 
		FROM place WHERE user_id = ? ORDER BY id DESC LIMIT ?`,
		userID, limit,
	)
	return items, err
}

func (r *Repository) GetPlaceByID(ctx context.Context, id, userID uint64) (*Place, error) {
	var p Place
	err := r.db.GetContext(ctx, &p,
		`SELECT id, user_id, name, link, link_type, description, go_at, go_at_time, status, created_at, updated_at 
		FROM place WHERE id = ? AND user_id = ?`,
		id, userID,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) UpdatePlace(ctx context.Context, id, userID uint64, req UpdatePlaceReq) (bool, error) {
	q := "UPDATE place SET "
	args := []any{}
	sets := []string{}

	if req.Name != nil {
		sets = append(sets, "name = ?")
		args = append(args, req.Name)
	}
	if req.Link != nil {
		sets = append(sets, "link = ?")
		args = append(args, req.Link)
	}
	if req.LinkType != nil {
		sets = append(sets, "link_type = ?")
		args = append(args, req.LinkType)
	}
	if req.Description != nil {
		sets = append(sets, "description = ?")
		args = append(args, req.Description)
	}
	if req.GoAt != nil {
		sets = append(sets, "go_at = ?")
		args = append(args, req.GoAt)
	}
	if req.GoAtTime != nil {
		sets = append(sets, "go_at_time = ?")
		args = append(args, req.GoAtTime)
	}
	if req.Status != nil {
		sets = append(sets, "status = ?")
		args = append(args, req.Status)
	}

	if len(sets) == 0 {
		return false, nil
	}

	q += strings.Join(sets, ", ") + " WHERE id = ? AND user_id = ?"
	args = append(args, id, userID)

	res, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (r *Repository) DeletePlace(ctx context.Context, id, userID uint64) (bool, error) {
	res, err := r.db.ExecContext(ctx, `DELETE FROM place WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

// PlaceCategory Repository Methods

func (r *Repository) CreatePlaceCategory(ctx context.Context, userID uint64, name string) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO place_category (user_id, name) VALUES (?, ?)`,
		userID, name,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) ListPlaceCategories(ctx context.Context, userID uint64, limit int) ([]PlaceCategory, error) {
	var items []PlaceCategory
	err := r.db.SelectContext(ctx, &items,
		`SELECT id, user_id, name FROM place_category WHERE user_id = ? ORDER BY id DESC LIMIT ?`,
		userID, limit,
	)
	return items, err
}

func (r *Repository) GetPlaceCategoryByID(ctx context.Context, id uint, userID uint64) (*PlaceCategory, error) {
	var pc PlaceCategory
	err := r.db.GetContext(ctx, &pc,
		`SELECT id, user_id, name FROM place_category WHERE id = ? AND user_id = ?`,
		id, userID,
	)
	if err != nil {
		return nil, err
	}
	return &pc, nil
}

func (r *Repository) UpdatePlaceCategory(ctx context.Context, id uint, userID uint64, name *string) (bool, error) {
	if name == nil {
		return false, nil
	}

	res, err := r.db.ExecContext(ctx,
		`UPDATE place_category SET name = ? WHERE id = ? AND user_id = ?`,
		*name, id, userID,
	)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (r *Repository) DeletePlaceCategory(ctx context.Context, id uint, userID uint64) (bool, error) {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM place_category WHERE id = ? AND user_id = ?`,
		id, userID,
	)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

// Helper function to convert Place model to response
func ToPlaceResponse(p *Place) PlaceResponse {
	resp := PlaceResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	}

	if p.Name.Valid {
		resp.Name = &p.Name.String
	}
	if p.Link.Valid {
		resp.Link = &p.Link.String
	}
	if p.LinkType.Valid {
		linkType := int(p.LinkType.Int32)
		resp.LinkType = &linkType
	}
	if p.Description.Valid {
		resp.Description = &p.Description.String
	}
	if p.GoAt.Valid {
		resp.GoAt = customtime.NewDate(p.GoAt.Time)
	}
	if p.GoAtTime.Valid {
		resp.GoAtTime = customtime.NewDateTime(p.GoAtTime.Time)
	}
	if p.Status.Valid {
		status := int(p.Status.Int32)
		resp.Status = &status
	}

	return resp
}

// Helper function to convert PlaceCategory model to response
func ToPlaceCategoryResponse(pc *PlaceCategory) PlaceCategoryResponse {
	return PlaceCategoryResponse{
		ID:     pc.ID,
		UserID: pc.UserID,
		Name:   pc.Name,
	}
}
