package user

import (
	"context"
	"net/http"
	"time"

	"go-saas-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service *Service
	v       *validator.Validate
}

func NewHandler(service *Service, v *validator.Validate) *Handler {
	return &Handler{
		service: service,
		v:       v,
	}
}

// POST /auth/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.v.Struct(req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	result, err := h.service.Register(ctx, req)
	if err != nil {
		if err.Error() == "email already registered" {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(c, http.StatusCreated, result)
}

// POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.v.Struct(req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	result, err := h.service.Login(ctx, req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			response.Error(c, http.StatusUnauthorized, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(c, http.StatusOK, result)
}

// POST /auth/change-password (protected route)
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.v.Struct(req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err := h.service.ChangePassword(ctx, userID.(uint64), req)
	if err != nil {
		if err.Error() == "old password is incorrect" {
			response.Error(c, http.StatusUnauthorized, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response.SuccessMessage(c, http.StatusOK, "password changed successfully")
}
