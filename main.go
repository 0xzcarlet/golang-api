package main

import (
	"context"
	"log"
	"os"
	"time"

	"go-saas-api/product"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is required")
	}

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("db ping failed:", err)
	}
	log.Println("db connected ✅")

	v := validator.New()
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	repo := product.NewRepo(db)
	handler := product.NewHandler(repo, v)
	product.RegisterRoutes(r, handler)

	log.Println("server running on :" + port)
	log.Fatal(r.Run(":" + port))
}
