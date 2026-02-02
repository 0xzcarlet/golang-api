package user

import (
	"go-saas-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, authMW *middleware.AuthMiddleware) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/change-password", authMW.RequireAuth(), h.ChangePassword)
	}
}
