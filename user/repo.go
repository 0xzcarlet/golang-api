package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

// Create - register user baru
func (r *Repo) Create(ctx context.Context, email, hashedPassword, name string) (uint64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users (email, password, name) VALUES (?, ?, ?)`,
		email, hashedPassword, name,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

// GetByEmail - untuk login
func (r *Repo) GetByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.GetContext(ctx, &u,
		`SELECT id, email, password, name, created_at, updated_at FROM users WHERE email = ?`,
		email,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByID - untuk verifikasi user dari token
func (r *Repo) GetByID(ctx context.Context, id uint64) (*User, error) {
	var u User
	err := r.db.GetContext(ctx, &u,
		`SELECT id, email, password, name, created_at, updated_at FROM users WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdatePassword - untuk change password
func (r *Repo) UpdatePassword(ctx context.Context, id uint64, hashedPassword string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password = ? WHERE id = ?`,
		hashedPassword, id,
	)
	return err
}

// EmailExists - cek apakah email sudah terdaftar
func (r *Repo) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM users WHERE email = ?`,
		email,
	)
	return count > 0, err
}
