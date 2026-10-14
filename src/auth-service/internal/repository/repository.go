package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"auth-service/internal/models"
)

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

type UserRepository interface {
	Create(ctx context.Context, u *models.User) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error
}
