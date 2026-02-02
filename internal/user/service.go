package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *Repository
	jwtSecret []byte
}

func NewService(repo *Repository, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *Service) Register(ctx context.Context, req RegisterReq) (*AuthResponse, error) {
	// Check if email already exists
	exists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	id, err := s.repo.Create(ctx, req.Email, string(hashedPassword), req.Name)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.generateToken(id)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    id,
			Email: req.Email,
			Name:  req.Name,
		},
	}, nil
}

func (s *Service) Login(ctx context.Context, req LoginReq) (*AuthResponse, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (s *Service) ChangePassword(ctx context.Context, userID uint64, req ChangePasswordReq) error {
	// Get current user
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	return s.repo.UpdatePassword(ctx, userID, string(hashedPassword))
}

func (s *Service) generateToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
