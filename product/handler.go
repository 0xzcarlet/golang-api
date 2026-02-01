package product

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	repo *Repo
	v    *validator.Validate
}

func NewHandler(repo *Repo, v *validator.Validate) *Handler {
	return &Handler{repo: repo, v: v}
}

// POST /products
func (h *Handler) Create(c *gin.Context) {
	var req CreateProductReq
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

	id, err := h.repo.Create(ctx, req.Name, req.Price)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "db error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})

}

// GET /products
func (h *Handler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	items, err := h.repo.List(ctx, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})

}

// GEt:id
func (h *Handler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrResp{Error: "invalid id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	p, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, ErrResp{Error: "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "db error"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// Update:id
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrResp{Error: "invalid id"})
		return
	}

	var req UpdateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrResp{Error: "invalid json"})
		return
	}

	// validasi: omitempty + min/max dll
	if err := h.v.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrResp{Error: "validation failed"})
		return
	}

	// body kosong: tidak ada field dikirim
	if req.Name == nil && req.Price == nil {
		c.JSON(http.StatusBadRequest, ErrResp{Error: "no fields to update"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Cek dulu apakah ID ada
	_, err = h.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, ErrResp{Error: "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "db error"})
		return
	}

	// Update data
	updated, err := h.repo.Update(ctx, id, req.Name, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "db error"})
		return
	}

	if !updated {
		c.JSON(http.StatusOK, gin.H{"message": "no changes made, data is already up to date"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update success"})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrResp{Error: "invalid id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	ok, err := h.repo.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrResp{Error: "db error"})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, ErrResp{Error: "not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
