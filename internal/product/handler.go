package product

import (
	"context"
	"net/http"
	"strconv"
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

// POST /products
func (h *Handler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateProductReq
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

	id, err := h.service.Create(ctx, userID.(uint64), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{"id": id})
}

// GET /products
func (h *Handler) List(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	items, err := h.service.List(ctx, userID.(uint64), 50)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"items": items})
}

// GET /products/:id
func (h *Handler) Get(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	product, err := h.service.GetByID(ctx, id, userID.(uint64))
	if err != nil {
		if err.Error() == "product not found" {
			response.Error(c, http.StatusNotFound, "product not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(c, http.StatusOK, product)
}

// PATCH /products/:id
func (h *Handler) Update(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateProductReq
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

	updated, err := h.service.Update(ctx, id, userID.(uint64), req)
	if err != nil {
		if err.Error() == "product not found" {
			response.Error(c, http.StatusNotFound, "product not found")
			return
		}
		if err.Error() == "no fields to update" {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if !updated {
		response.SuccessMessage(c, http.StatusOK, "no changes made, data is already up to date")
		return
	}

	response.SuccessMessage(c, http.StatusOK, "product updated successfully")
}

// DELETE /products/:id
func (h *Handler) Delete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	_, err = h.service.Delete(ctx, id, userID.(uint64))
	if err != nil {
		if err.Error() == "product not found" {
			response.Error(c, http.StatusNotFound, "product not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.Status(http.StatusNoContent)
}
