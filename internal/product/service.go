package product

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

func (s *Service) Create(ctx context.Context, userID uint64, req CreateProductReq) (int64, error) {
	return s.repo.Create(ctx, userID, req.Name, req.Price)
}

func (s *Service) List(ctx context.Context, userID uint64, limit int) ([]Product, error) {
	return s.repo.List(ctx, userID, limit)
}

func (s *Service) GetByID(ctx context.Context, id, userID uint64) (*Product, error) {
	product, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return product, nil
}

func (s *Service) Update(ctx context.Context, id, userID uint64, req UpdateProductReq) (bool, error) {
	if req.Name == nil && req.Price == nil {
		return false, errors.New("no fields to update")
	}

	// Check if product exists and belongs to user
	_, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("product not found")
		}
		return false, err
	}

	return s.repo.Update(ctx, id, userID, req.Name, req.Price)
}

func (s *Service) Delete(ctx context.Context, id, userID uint64) (bool, error) {
	deleted, err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		return false, err
	}
	if !deleted {
		return false, errors.New("product not found")
	}
	return true, nil
}
