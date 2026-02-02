package place

import (
	"context"
	"database/sql"
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Place Service Methods

func (s *Service) CreatePlace(ctx context.Context, userID uint64, req CreatePlaceReq) (int64, error) {
	return s.repo.CreatePlace(ctx, userID, req)
}

func (s *Service) ListPlaces(ctx context.Context, userID uint64, limit int) ([]Place, error) {
	return s.repo.ListPlaces(ctx, userID, limit)
}

func (s *Service) GetPlaceByID(ctx context.Context, id, userID uint64) (*Place, error) {
	place, err := s.repo.GetPlaceByID(ctx, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("place not found")
		}
		return nil, err
	}
	return place, nil
}

func (s *Service) UpdatePlace(ctx context.Context, id, userID uint64, req UpdatePlaceReq) (bool, error) {
	if req.Name == nil && req.Link == nil && req.LinkType == nil &&
		req.Description == nil && req.GoAt == nil && req.GoAtTime == nil && req.Status == nil {
		return false, errors.New("no fields to update")
	}

	// Check if place exists
	_, err := s.repo.GetPlaceByID(ctx, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("place not found")
		}
		return false, err
	}

	return s.repo.UpdatePlace(ctx, id, userID, req)
}

func (s *Service) DeletePlace(ctx context.Context, id, userID uint64) (bool, error) {
	deleted, err := s.repo.DeletePlace(ctx, id, userID)
	if err != nil {
		return false, err
	}
	if !deleted {
		return false, errors.New("place not found")
	}
	return true, nil
}

// PlaceCategory Service Methods

func (s *Service) CreatePlaceCategory(ctx context.Context, userID uint64, req CreatePlaceCategoryReq) (int64, error) {
	return s.repo.CreatePlaceCategory(ctx, userID, req.Name)
}

func (s *Service) ListPlaceCategories(ctx context.Context, userID uint64, limit int) ([]PlaceCategory, error) {
	return s.repo.ListPlaceCategories(ctx, userID, limit)
}

func (s *Service) GetPlaceCategoryByID(ctx context.Context, id uint, userID uint64) (*PlaceCategory, error) {
	category, err := s.repo.GetPlaceCategoryByID(ctx, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return category, nil
}

func (s *Service) UpdatePlaceCategory(ctx context.Context, id uint, userID uint64, req UpdatePlaceCategoryReq) (bool, error) {
	if req.Name == nil {
		return false, errors.New("no fields to update")
	}

	// Check if category exists and belongs to user
	_, err := s.repo.GetPlaceCategoryByID(ctx, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("category not found")
		}
		return false, err
	}

	return s.repo.UpdatePlaceCategory(ctx, id, userID, req.Name)
}

func (s *Service) DeletePlaceCategory(ctx context.Context, id uint, userID uint64) (bool, error) {
	deleted, err := s.repo.DeletePlaceCategory(ctx, id, userID)
	if err != nil {
		return false, err
	}
	if !deleted {
		return false, errors.New("category not found")
	}
	return true, nil
}
