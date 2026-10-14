package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"shop-api/internal/models"
)

type clientRepo struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) ClientRepository {
	return &clientRepo{db: db}
}

func (r *clientRepo) Create(ctx context.Context, c *models.Client, a *models.Address) (*models.Client, *models.Address, error) {
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

	out := models.Client{
		ClientName:    c.ClientName,
		ClientSurname: c.ClientSurname,
		Birthday:      c.Birthday,
		Gender:        c.Gender,
		AddressID:     a.ID,
	}
	if err = tx.QueryRowContext(ctx,
		`INSERT INTO clients (client_name, client_surname, birthday, gender, address_id)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, registration_date`,
		out.ClientName, out.ClientSurname, out.Birthday, out.Gender, out.AddressID,
	).Scan(&out.ID, &out.RegistrationDate); err != nil {
		return nil, nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, nil, err
	}
	return &out, a, nil
}

func (r *clientRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Client, *models.Address, error) {
	var c models.Client
	var a models.Address
	err := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.client_name, c.client_surname, c.birthday, c.gender,
		       c.registration_date, c.address_id,
		       a.id, a.country, a.city, a.street
		FROM clients c
		JOIN addresses a ON a.id = c.address_id
		WHERE c.id = $1`, id,
	).Scan(
		&c.ID, &c.ClientName, &c.ClientSurname, &c.Birthday, &c.Gender,
		&c.RegistrationDate, &c.AddressID,
		&a.ID, &a.Country, &a.City, &a.Street,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, ErrNotFound
	}
	if err != nil {
		return nil, nil, err
	}
	return &c, &a, nil
}

func (r *clientRepo) List(ctx context.Context, name, surname string, limit, offset int) ([]models.ClientFull, error) {
	q := `
		SELECT c.id, c.client_name, c.client_surname, c.birthday, c.gender,
		       c.registration_date, c.address_id,
		       a.id, a.country, a.city, a.street
		FROM clients c
		JOIN addresses a ON a.id = c.address_id
		WHERE 1=1`
	args := []any{}

	if name != "" {
		args = append(args, name)
		q += fmt.Sprintf(" AND c.client_name = $%d", len(args))
	}
	if surname != "" {
		args = append(args, surname)
		q += fmt.Sprintf(" AND c.client_surname = $%d", len(args))
	}
	q += " ORDER BY c.registration_date DESC"

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

	out := make([]models.ClientFull, 0)
	for rows.Next() {
		var cf models.ClientFull
		if err := rows.Scan(
			&cf.ID, &cf.ClientName, &cf.ClientSurname, &cf.Birthday, &cf.Gender,
			&cf.RegistrationDate, &cf.AddressID,
			&cf.Address.ID, &cf.Address.Country, &cf.Address.City, &cf.Address.Street,
		); err != nil {
			return nil, err
		}
		out = append(out, cf)
	}
	return out, rows.Err()
}

func (r *clientRepo) UpdateAddress(ctx context.Context, clientID uuid.UUID, a *models.Address) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE addresses SET country=$1, city=$2, street=$3
		WHERE id = (SELECT address_id FROM clients WHERE id=$4)`,
		a.Country, a.City, a.Street, clientID,
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

func (r *clientRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var addressID uuid.UUID
	if err = tx.QueryRowContext(ctx, `SELECT address_id FROM clients WHERE id=$1`, id).Scan(&addressID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM clients WHERE id=$1`, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM addresses WHERE id=$1`, addressID); err != nil {
		return err
	}
	return tx.Commit()
}
