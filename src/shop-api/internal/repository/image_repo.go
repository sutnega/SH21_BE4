package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"shop-api/internal/models"
)

type imageRepo struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) ImageRepository {
	return &imageRepo{db: db}
}

func (r *imageRepo) Create(ctx context.Context, productID uuid.UUID, data []byte) (*models.Image, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var img models.Image
	if err = tx.QueryRowContext(ctx,
		`INSERT INTO images (image) VALUES ($1) RETURNING id`, data,
	).Scan(&img.ID); err != nil {
		return nil, err
	}
	img.Image = data

	res, err := tx.ExecContext(ctx, `UPDATE products SET image_id=$1 WHERE id=$2`, img.ID, productID)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, ErrNotFound
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *imageRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Image, error) {
	var img models.Image
	err := r.db.QueryRowContext(ctx, `SELECT id, image FROM images WHERE id=$1`, id).
		Scan(&img.ID, &img.Image)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *imageRepo) GetByProductID(ctx context.Context, productID uuid.UUID) (*models.Image, error) {
	var img models.Image
	err := r.db.QueryRowContext(ctx, `
		SELECT i.id, i.image
		FROM images i
		JOIN products p ON p.image_id = i.id
		WHERE p.id = $1`, productID,
	).Scan(&img.ID, &img.Image)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *imageRepo) Update(ctx context.Context, id uuid.UUID, data []byte) error {
	res, err := r.db.ExecContext(ctx, `UPDATE images SET image=$1 WHERE id=$2`, data, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *imageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM images WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
