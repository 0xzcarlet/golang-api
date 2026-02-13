package main

import (
	"log"

	"go-saas-api/internal/config"
	"go-saas-api/internal/database"
	"go-saas-api/internal/middleware"
	"go-saas-api/internal/place"
	"go-saas-api/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup database connection
	db, err := database.NewConnection(cfg.DBDsn)
	if err != nil {
		log.Fatal("database connection failed:", err)
	}
	defer db.Close()
	log.Println("✅ Database connected")

	// Setup validator
	v := validator.New()

	// Setup middleware
	authMW := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Setup Gin router
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Setup modules
	setupUserModule(r, db, v, cfg.JWTSecret, authMW)
	setupPlaceModule(r, db, v, authMW)

	// Start server
	log.Printf("🚀 Server running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("server failed to start:", err)
	}
}

func setupUserModule(r *gin.Engine, db *sqlx.DB, v *validator.Validate, jwtSecret string, authMW *middleware.AuthMiddleware) {
	repo := user.NewRepository(db)
	service := user.NewService(repo, jwtSecret)
	handler := user.NewHandler(service, v)
	user.RegisterRoutes(r, handler, authMW)
}

func setupPlaceModule(r *gin.Engine, db *sqlx.DB, v *validator.Validate, authMW *middleware.AuthMiddleware) {
	repo := place.NewRepository(db)
	service := place.NewService(repo)
	handler := place.NewHandler(service, v)
	place.RegisterRoutes(r, handler, authMW)
}
