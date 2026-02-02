package product

import (
	"go-saas-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, authMW *middleware.AuthMiddleware) {
	products := r.Group("/products", authMW.RequireAuth())
	{
		products.POST("", h.Create)
		products.GET("", h.List)
		products.GET("/:id", h.Get)
		products.PATCH("/:id", h.Update)
		products.DELETE("/:id", h.Delete)
	}
}
