package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewConnection(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
