package product

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

func (r *Repository) Create(ctx context.Context, userID uint64, name string, price int) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO products (user_id, name, price) VALUES (?, ?, ?)`,
		userID, name, price,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) List(ctx context.Context, userID uint64, limit int) ([]Product, error) {
	var items []Product
	err := r.db.SelectContext(ctx, &items,
		`SELECT id, user_id, name, price, created_at, updated_at FROM products WHERE user_id = ? ORDER BY id DESC LIMIT ?`,
		userID, limit,
	)
	return items, err
}

func (r *Repository) GetByID(ctx context.Context, id, userID uint64) (*Product, error) {
	var p Product
	err := r.db.GetContext(ctx, &p,
		`SELECT id, user_id, name, price, created_at, updated_at FROM products WHERE id = ? AND user_id = ?`,
		id, userID,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) Update(ctx context.Context, id, userID uint64, name *string, price *int) (bool, error) {
	q := "UPDATE products SET "
	args := []any{}
	set := ""

	if name != nil {
		set += "name = ?"
		args = append(args, *name)
	}
	if price != nil {
		if set != "" {
			set += ", "
		}
		set += "price = ?"
		args = append(args, *price)
	}

	if set == "" {
		return false, nil
	}

	q += set + " WHERE id = ? AND user_id = ?"
	args = append(args, id, userID)

	res, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (r *Repository) Delete(ctx context.Context, id, userID uint64) (bool, error) {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM products WHERE id = ? AND user_id = ?`,
		id, userID,
	)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return true, nil
	}
	return rows > 0, nil
}
