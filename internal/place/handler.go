package place

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

// Place Handlers

// POST /places
func (h *Handler) CreatePlace(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreatePlaceReq
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

	id, err := h.service.CreatePlace(ctx, userID.(uint64), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{"id": id})
}

// GET /places
func (h *Handler) ListPlaces(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	items, err := h.service.ListPlaces(ctx, userID.(uint64), 100)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	// Convert to response format
	responses := make([]PlaceResponse, len(items))
	for i, item := range items {
		responses[i] = ToPlaceResponse(&item)
	}

	response.Success(c, http.StatusOK, gin.H{"items": responses})
}

// GET /places/:id
func (h *Handler) GetPlace(c *gin.Context) {
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

	place, err := h.service.GetPlaceByID(ctx, id, userID.(uint64))
	if err != nil {
		if err.Error() == "place not found" {
			response.Error(c, http.StatusNotFound, "place not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := ToPlaceResponse(place)
	response.Success(c, http.StatusOK, resp)
}

// PATCH /places/:id
func (h *Handler) UpdatePlace(c *gin.Context) {
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

	var req UpdatePlaceReq
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

	updated, err := h.service.UpdatePlace(ctx, id, userID.(uint64), req)
	if err != nil {
		if err.Error() == "place not found" {
			response.Error(c, http.StatusNotFound, "place not found")
			return
		}
		if err.Error() == "no fields to update" {
			response.Error(c, http.StatusBadRequest, "no fields to update")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if !updated {
		response.Error(c, http.StatusNotFound, "place not found")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "place updated successfully"})
}

// DELETE /places/:id
func (h *Handler) DeletePlace(c *gin.Context) {
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

	deleted, err := h.service.DeletePlace(ctx, id, userID.(uint64))
	if err != nil {
		if err.Error() == "place not found" {
			response.Error(c, http.StatusNotFound, "place not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if !deleted {
		response.Error(c, http.StatusNotFound, "place not found")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "place deleted successfully"})
}

// PlaceCategory Handlers

// POST /place-categories
func (h *Handler) CreatePlaceCategory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreatePlaceCategoryReq
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

	id, err := h.service.CreatePlaceCategory(ctx, userID.(uint64), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{"id": id})
}

// GET /place-categories
func (h *Handler) ListPlaceCategories(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	items, err := h.service.ListPlaceCategories(ctx, userID.(uint64), 100)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	// Convert to response format
	responses := make([]PlaceCategoryResponse, len(items))
	for i, item := range items {
		responses[i] = ToPlaceCategoryResponse(&item)
	}

	response.Success(c, http.StatusOK, gin.H{"items": responses})
}

// GET /place-categories/:id
func (h *Handler) GetPlaceCategory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	category, err := h.service.GetPlaceCategoryByID(ctx, uint(id), userID.(uint64))
	if err != nil {
		if err.Error() == "category not found" {
			response.Error(c, http.StatusNotFound, "category not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := ToPlaceCategoryResponse(category)
	response.Success(c, http.StatusOK, resp)
}

// PATCH /place-categories/:id
func (h *Handler) UpdatePlaceCategory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdatePlaceCategoryReq
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

	updated, err := h.service.UpdatePlaceCategory(ctx, uint(id), userID.(uint64), req)
	if err != nil {
		if err.Error() == "category not found" {
			response.Error(c, http.StatusNotFound, "category not found")
			return
		}
		if err.Error() == "no fields to update" {
			response.Error(c, http.StatusBadRequest, "no fields to update")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if !updated {
		response.Error(c, http.StatusNotFound, "category not found")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "category updated successfully"})
}

// DELETE /place-categories/:id
func (h *Handler) DeletePlaceCategory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	deleted, err := h.service.DeletePlaceCategory(ctx, uint(id), userID.(uint64))
	if err != nil {
		if err.Error() == "category not found" {
			response.Error(c, http.StatusNotFound, "category not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if !deleted {
		response.Error(c, http.StatusNotFound, "category not found")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "category deleted successfully"})
}
