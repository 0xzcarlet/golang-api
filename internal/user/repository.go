package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, email, hashedPassword, name string) (uint64, error) {
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

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
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

func (r *Repository) GetByID(ctx context.Context, id uint64) (*User, error) {
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

func (r *Repository) UpdatePassword(ctx context.Context, id uint64, hashedPassword string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password = ? WHERE id = ?`,
		hashedPassword, id,
	)
	return err
}

func (r *Repository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM users WHERE email = ?`,
		email,
	)
	return count > 0, err
}
