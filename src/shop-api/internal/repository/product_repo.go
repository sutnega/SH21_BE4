package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"shop-api/internal/models"
)

type productRepo struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepo{db: db}
}

const productSelect = `
	SELECT p.id, p.name, p.category_id, p.price, p.available_stock,
	       p.last_update_date, p.supplier_id, p.image_id,
	       c.name
	FROM products p
	JOIN categories c ON c.id = p.category_id
`

func scanProductFull(rows interface {
	Scan(dest ...any) error
}) (models.ProductFull, error) {
	var p models.ProductFull
	err := rows.Scan(
		&p.ID, &p.Name, &p.CategoryID, &p.Price, &p.AvailableStock,
		&p.LastUpdateDate, &p.SupplierID, &p.ImageID,
		&p.CategoryName,
	)
	return p, err
}

func (r *productRepo) Create(ctx context.Context, p *models.Product, categoryName string) (*models.Product, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var categoryID uuid.UUID
	err = tx.QueryRowContext(ctx,
		`SELECT id FROM categories WHERE name=$1`, categoryName,
	).Scan(&categoryID)
	if errors.Is(err, sql.ErrNoRows) {
		err = tx.QueryRowContext(ctx,
			`INSERT INTO categories(name) VALUES($1) RETURNING id`, categoryName,
		).Scan(&categoryID)
	}
	if err != nil {
		return nil, err
	}

	out := models.Product{
		Name:           p.Name,
		CategoryID:     categoryID,
		Price:          p.Price,
		AvailableStock: p.AvailableStock,
		SupplierID:     p.SupplierID,
	}
	if err = tx.QueryRowContext(ctx, `
		INSERT INTO products (name, category_id, price, available_stock, supplier_id)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, last_update_date, image_id`,
		out.Name, out.CategoryID, out.Price, out.AvailableStock, out.SupplierID,
	).Scan(&out.ID, &out.LastUpdateDate, &out.ImageID); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *productRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.ProductFull, error) {
	row := r.db.QueryRowContext(ctx, productSelect+` WHERE p.id = $1`, id)
	p, err := scanProductFull(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepo) ListAvailable(ctx context.Context, limit, offset int) ([]models.ProductFull, error) {
	q := productSelect + ` WHERE p.available_stock > 0 ORDER BY p.last_update_date DESC`
	args := []any{}
	if limit > 0 {
		args = append(args, limit)
		q += fmt.Sprintf(" LIMIT $%d", len(args))
	}
	if offset > 0 {
		args = append(args, offset)
		q += fmt.Sprintf(" OFFSET $%d", len(args))
	}

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.ProductFull, 0)
	for rows.Next() {
		p, err := scanProductFull(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *productRepo) DecreaseStock(ctx context.Context, id uuid.UUID, amount int) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE products
		SET available_stock = available_stock - $1,
		    last_update_date = NOW()
		WHERE id = $2 AND available_stock >= $1`,
		amount, id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		// либо нет товара, либо недостаточно
		var exists bool
		if err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)`, id).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return ErrNotFound
		}
		return errors.New("not enough stock")
	}
	return nil
}

func (r *productRepo) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
