package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"shop-api/internal/models"
)

type supplierRepo struct {
	db *sql.DB
}

func NewSupplierRepository(db *sql.DB) SupplierRepository {
	return &supplierRepo{db: db}
}

func (r *supplierRepo) Create(ctx context.Context, s *models.Supplier, a *models.Address) (*models.Supplier, *models.Address, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	if err = tx.QueryRowContext(ctx,
		`INSERT INTO addresses (country, city, street) VALUES ($1,$2,$3) RETURNING id`,
		a.Country, a.City, a.Street,
	).Scan(&a.ID); err != nil {
		return nil, nil, err
	}

	out := models.Supplier{Name: s.Name, AddressID: a.ID, PhoneNumber: s.PhoneNumber}
	if err = tx.QueryRowContext(ctx,
		`INSERT INTO suppliers (name, address_id, phone_number)
		 VALUES ($1,$2,$3) RETURNING id`,
		out.Name, out.AddressID, out.PhoneNumber,
	).Scan(&out.ID); err != nil {
		return nil, nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, nil, err
	}
	return &out, a, nil
}

func (r *supplierRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Supplier, *models.Address, error) {
	var s models.Supplier
	var a models.Address
	err := r.db.QueryRowContext(ctx, `
		SELECT s.id, s.name, s.address_id, s.phone_number,
		       a.id, a.country, a.city, a.street
		FROM suppliers s
		JOIN addresses a ON a.id = s.address_id
		WHERE s.id = $1`, id,
	).Scan(
		&s.ID, &s.Name, &s.AddressID, &s.PhoneNumber,
		&a.ID, &a.Country, &a.City, &a.Street,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, ErrNotFound
	}
	if err != nil {
		return nil, nil, err
	}
	return &s, &a, nil
}

func (r *supplierRepo) List(ctx context.Context) ([]models.SupplierFull, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.name, s.address_id, s.phone_number,
		       a.id, a.country, a.city, a.street
		FROM suppliers s
		JOIN addresses a ON a.id = s.address_id
		ORDER BY s.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.SupplierFull, 0)
	for rows.Next() {
		var sf models.SupplierFull
		if err := rows.Scan(
			&sf.ID, &sf.Name, &sf.AddressID, &sf.PhoneNumber,
			&sf.Address.ID, &sf.Address.Country, &sf.Address.City, &sf.Address.Street,
		); err != nil {
			return nil, err
		}
		out = append(out, sf)
	}
	return out, rows.Err()
}

func (r *supplierRepo) UpdateAddress(ctx context.Context, supplierID uuid.UUID, a *models.Address) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE addresses SET country=$1, city=$2, street=$3
		WHERE id = (SELECT address_id FROM suppliers WHERE id=$4)`,
		a.Country, a.City, a.Street, supplierID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *supplierRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var addressID uuid.UUID
	if err = tx.QueryRowContext(ctx, `SELECT address_id FROM suppliers WHERE id=$1`, id).Scan(&addressID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM suppliers WHERE id=$1`, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM addresses WHERE id=$1`, addressID); err != nil {
		return err
	}
	return tx.Commit()
}
