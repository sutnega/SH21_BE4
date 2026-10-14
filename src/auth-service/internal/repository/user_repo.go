package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"auth-service/internal/models"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *models.User) (*models.User, error) {
	out := models.User{
		Email:        u.Email,
		Name:         u.Name,
		Surname:      u.Surname,
		Phone:        u.Phone,
		PasswordHash: u.PasswordHash,
	}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (email, name, surname, phone, password_hash)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, created_at`,
		out.Email, out.Name, out.Surname, out.Phone, out.PasswordHash,
	).Scan(&out.ID, &out.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, ErrAlreadyExists
		}
		return nil, err
	}
	return &out, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, name, surname, phone, password_hash, created_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.Name, &u.Surname, &u.Phone, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, name, surname, phone, password_hash, created_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.Surname, &u.Phone, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash = $1 WHERE id = $2`, hash, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
