package user

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo      *Repo
	v         *validator.Validate
	jwtSecret []byte
}

func NewHandler(repo *Repo, v *validator.Validate) *Handler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production" // default untuk development
	}
	return &Handler{
		repo:      repo,
		v:         v,
		jwtSecret: []byte(secret),
	}
}

// POST /auth/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrResp{Error: "invalid json"})
		return
	}
	if err := h.v.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrResp{Error: "validation failed"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Cek apakah email sudah ada
	exists, err := h.repo.EmailExists(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "db error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, ErrResp{Error: "email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "internal error"})
		return
	}

	// Create user
	id, err := h.repo.Create(ctx, req.Email, string(hashedPassword), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "db error"})
		return
	}

	// Generate JWT token
	token, err := h.generateToken(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "token generation failed"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    id,
			Email: req.Email,
			Name:  req.Name,
		},
	})
}

// POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrResp{Error: "invalid json"})
		return
	}
	if err := h.v.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrResp{Error: "validation failed"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get user by email
	user, err := h.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, ErrResp{Error: "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "db error"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrResp{Error: "invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	})
}

// POST /auth/change-password (protected route)
func (h *Handler) ChangePassword(c *gin.Context) {
	// Get user ID from context (set by middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrResp{Error: "unauthorized"})
		return
	}

	var req ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrResp{Error: "invalid json"})
		return
	}
	if err := h.v.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrResp{Error: "validation failed"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get current user
	user, err := h.repo.GetByID(ctx, userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "db error"})
		return
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrResp{Error: "old password is incorrect"})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "internal error"})
		return
	}

	// Update password
	err = h.repo.UpdatePassword(ctx, userID.(uint64), string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// Helper: generate JWT token
func (h *Handler) generateToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // expired 24 jam
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.jwtSecret)
}
