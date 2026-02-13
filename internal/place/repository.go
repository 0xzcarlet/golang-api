package place

import (
	"context"
	"fmt"
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
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO place (user_id, name, link, link_type, description, go_at, go_at_time, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		userID, req.Name, req.Link, req.LinkType, req.Description, req.GoAt, req.GoAtTime, req.Status,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) ListPlaces(ctx context.Context, userID uint64, limit int) ([]Place, error) {
	var items []Place
	err := r.db.SelectContext(ctx, &items,
		`SELECT id, user_id, name, link, link_type, description, go_at, go_at_time, status, created_at, updated_at 
		FROM place WHERE user_id = $1 ORDER BY id DESC LIMIT $2`,
		userID, limit,
	)
	return items, err
}

func (r *Repository) GetPlaceByID(ctx context.Context, id, userID uint64) (*Place, error) {
	var p Place
	err := r.db.GetContext(ctx, &p,
		`SELECT id, user_id, name, link, link_type, description, go_at, go_at_time, status, created_at, updated_at 
		FROM place WHERE id = $1 AND user_id = $2`,
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
	paramIdx := 1

	if req.Name != nil {
		sets = append(sets, fmt.Sprintf("name = $%d", paramIdx))
		args = append(args, req.Name)
		paramIdx++
	}
	if req.Link != nil {
		sets = append(sets, fmt.Sprintf("link = $%d", paramIdx))
		args = append(args, req.Link)
		paramIdx++
	}
	if req.LinkType != nil {
		sets = append(sets, fmt.Sprintf("link_type = $%d", paramIdx))
		args = append(args, req.LinkType)
		paramIdx++
	}
	if req.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", paramIdx))
		args = append(args, req.Description)
		paramIdx++
	}
	if req.GoAt != nil {
		sets = append(sets, fmt.Sprintf("go_at = $%d", paramIdx))
		args = append(args, req.GoAt)
		paramIdx++
	}
	if req.GoAtTime != nil {
		sets = append(sets, fmt.Sprintf("go_at_time = $%d", paramIdx))
		args = append(args, req.GoAtTime)
		paramIdx++
	}
	if req.Status != nil {
		sets = append(sets, fmt.Sprintf("status = $%d", paramIdx))
		args = append(args, req.Status)
		paramIdx++
	}

	if len(sets) == 0 {
		return false, nil
	}

	q += strings.Join(sets, ", ") + fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", paramIdx, paramIdx+1)
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
	res, err := r.db.ExecContext(ctx, `DELETE FROM place WHERE id = $1 AND user_id = $2`, id, userID)
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
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO place_category (user_id, name) VALUES ($1, $2) RETURNING id`,
		userID, name,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) ListPlaceCategories(ctx context.Context, userID uint64, limit int) ([]PlaceCategory, error) {
	var items []PlaceCategory
	err := r.db.SelectContext(ctx, &items,
		`SELECT id, user_id, name FROM place_category WHERE user_id = $1 ORDER BY id DESC LIMIT $2`,
		userID, limit,
	)
	return items, err
}

func (r *Repository) GetPlaceCategoryByID(ctx context.Context, id uint, userID uint64) (*PlaceCategory, error) {
	var pc PlaceCategory
	err := r.db.GetContext(ctx, &pc,
		`SELECT id, user_id, name FROM place_category WHERE id = $1 AND user_id = $2`,
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
		`UPDATE place_category SET name = $1 WHERE id = $2 AND user_id = $3`,
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
		`DELETE FROM place_category WHERE id = $1 AND user_id = $2`,
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
