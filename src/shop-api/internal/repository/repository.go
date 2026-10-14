package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"shop-api/internal/models"
)

var ErrNotFound = errors.New("not found")

type ClientRepository interface {
	Create(ctx context.Context, c *models.Client, a *models.Address) (*models.Client, *models.Address, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Client, *models.Address, error)
	List(ctx context.Context, name, surname string, limit, offset int) ([]models.ClientFull, error)
	UpdateAddress(ctx context.Context, clientID uuid.UUID, a *models.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProductRepository interface {
	Create(ctx context.Context, p *models.Product, categoryName string) (*models.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.ProductFull, error)
	ListAvailable(ctx context.Context, limit, offset int) ([]models.ProductFull, error)
	DecreaseStock(ctx context.Context, id uuid.UUID, amount int) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SupplierRepository interface {
	Create(ctx context.Context, s *models.Supplier, a *models.Address) (*models.Supplier, *models.Address, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Supplier, *models.Address, error)
	List(ctx context.Context) ([]models.SupplierFull, error)
	UpdateAddress(ctx context.Context, supplierID uuid.UUID, a *models.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ImageRepository interface {
	Create(ctx context.Context, productID uuid.UUID, data []byte) (*models.Image, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Image, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) (*models.Image, error)
	Update(ctx context.Context, id uuid.UUID, data []byte) error
	Delete(ctx context.Context, id uuid.UUID) error
}
