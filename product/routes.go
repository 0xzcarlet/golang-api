package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.POST("/products", h.Create)
	r.GET("/products", h.List)
	r.GET("/products/:id", h.Get)
	r.PATCH("/products/:id", h.Update)
	r.DELETE("products/:id", h.Delete)
}
