package place

import (
	"go-saas-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, authMW *middleware.AuthMiddleware) {
	// Place routes - require authentication
	places := r.Group("/places", authMW.RequireAuth())
	{
		places.POST("", h.CreatePlace)
		places.GET("", h.ListPlaces)
		places.GET("/:id", h.GetPlace)
		places.PATCH("/:id", h.UpdatePlace)
		places.DELETE("/:id", h.DeletePlace)
	}

	// PlaceCategory routes - require authentication
	categories := r.Group("/place-categories", authMW.RequireAuth())
	{
		categories.POST("", h.CreatePlaceCategory)
		categories.GET("", h.ListPlaceCategories)
		categories.GET("/:id", h.GetPlaceCategory)
		categories.PATCH("/:id", h.UpdatePlaceCategory)
		categories.DELETE("/:id", h.DeletePlaceCategory)
	}
}
