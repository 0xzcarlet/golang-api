package product

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

func (r *Repo) Create(ctx context.Context, name string, price int) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO products (name, price) VALUES (?, ?)`,
		name, price,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repo) List(ctx context.Context, limit int) ([]Product, error) {
	var items []Product
	err := r.db.SelectContext(ctx, &items,
		`SELECT id, name, price, created_at FROM products ORDER BY id DESC LIMIT ?`,
		limit,
	)
	return items, err
}

func (r *Repo) GetByID(ctx context.Context, id uint64) (*Product, error) {
	var p Product
	err := r.db.GetContext(ctx, &p,
		`SELECT id, name, price, created_at FROM products WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repo) Update(ctx context.Context, id uint64, name *string, price *int) (bool, error) {
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
		// tidak ada field yang mau diupdate
		return false, nil
	}

	q += set + " WHERE id = ?"
	args = append(args, id)

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

func (r *Repo) Delete(ctx context.Context, id uint64) (bool, error) {
	res, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return true, nil
	}
	return rows > 0, nil
}
