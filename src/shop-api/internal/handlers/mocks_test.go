package handlers

import (
	"context"

	"github.com/google/uuid"

	"shop-api/internal/models"
	"shop-api/internal/repository"
)

// ---------- Client ----------

type mockClientRepo struct {
	createFn        func(ctx context.Context, c *models.Client, a *models.Address) (*models.Client, *models.Address, error)
	getByIDFn       func(ctx context.Context, id uuid.UUID) (*models.Client, *models.Address, error)
	listFn          func(ctx context.Context, name, surname string, limit, offset int) ([]models.ClientFull, error)
	updateAddressFn func(ctx context.Context, id uuid.UUID, a *models.Address) error
	deleteFn        func(ctx context.Context, id uuid.UUID) error
}

func (m *mockClientRepo) Create(ctx context.Context, c *models.Client, a *models.Address) (*models.Client, *models.Address, error) {
	return m.createFn(ctx, c, a)
}
func (m *mockClientRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Client, *models.Address, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockClientRepo) List(ctx context.Context, name, surname string, limit, offset int) ([]models.ClientFull, error) {
	return m.listFn(ctx, name, surname, limit, offset)
}
func (m *mockClientRepo) UpdateAddress(ctx context.Context, id uuid.UUID, a *models.Address) error {
	return m.updateAddressFn(ctx, id, a)
}
func (m *mockClientRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}

// ---------- Product ----------

type mockProductRepo struct {
	createFn        func(ctx context.Context, p *models.Product, categoryName string) (*models.Product, error)
	getByIDFn       func(ctx context.Context, id uuid.UUID) (*models.ProductFull, error)
	listAvailableFn func(ctx context.Context, limit, offset int) ([]models.ProductFull, error)
	decreaseStockFn func(ctx context.Context, id uuid.UUID, amount int) error
	deleteFn        func(ctx context.Context, id uuid.UUID) error
}

func (m *mockProductRepo) Create(ctx context.Context, p *models.Product, c string) (*models.Product, error) {
	return m.createFn(ctx, p, c)
}
func (m *mockProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.ProductFull, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockProductRepo) ListAvailable(ctx context.Context, limit, offset int) ([]models.ProductFull, error) {
	return m.listAvailableFn(ctx, limit, offset)
}
func (m *mockProductRepo) DecreaseStock(ctx context.Context, id uuid.UUID, amount int) error {
	return m.decreaseStockFn(ctx, id, amount)
}
func (m *mockProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}

// ---------- Supplier ----------

type mockSupplierRepo struct {
	createFn        func(ctx context.Context, s *models.Supplier, a *models.Address) (*models.Supplier, *models.Address, error)
	getByIDFn       func(ctx context.Context, id uuid.UUID) (*models.Supplier, *models.Address, error)
	listFn          func(ctx context.Context) ([]models.SupplierFull, error)
	updateAddressFn func(ctx context.Context, id uuid.UUID, a *models.Address) error
	deleteFn        func(ctx context.Context, id uuid.UUID) error
}

func (m *mockSupplierRepo) Create(ctx context.Context, s *models.Supplier, a *models.Address) (*models.Supplier, *models.Address, error) {
	return m.createFn(ctx, s, a)
}
func (m *mockSupplierRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Supplier, *models.Address, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockSupplierRepo) List(ctx context.Context) ([]models.SupplierFull, error) {
	return m.listFn(ctx)
}
func (m *mockSupplierRepo) UpdateAddress(ctx context.Context, id uuid.UUID, a *models.Address) error {
	return m.updateAddressFn(ctx, id, a)
}
func (m *mockSupplierRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}

// ---------- Image ----------

type mockImageRepo struct {
	createFn         func(ctx context.Context, productID uuid.UUID, data []byte) (*models.Image, error)
	getByIDFn        func(ctx context.Context, id uuid.UUID) (*models.Image, error)
	getByProductIDFn func(ctx context.Context, productID uuid.UUID) (*models.Image, error)
	updateFn         func(ctx context.Context, id uuid.UUID, data []byte) error
	deleteFn         func(ctx context.Context, id uuid.UUID) error
}

func (m *mockImageRepo) Create(ctx context.Context, productID uuid.UUID, data []byte) (*models.Image, error) {
	return m.createFn(ctx, productID, data)
}
func (m *mockImageRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Image, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockImageRepo) GetByProductID(ctx context.Context, productID uuid.UUID) (*models.Image, error) {
	return m.getByProductIDFn(ctx, productID)
}
func (m *mockImageRepo) Update(ctx context.Context, id uuid.UUID, data []byte) error {
	return m.updateFn(ctx, id, data)
}
func (m *mockImageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}

// проверка, что моки удовлетворяют интерфейсам
var (
	_ repository.ClientRepository   = (*mockClientRepo)(nil)
	_ repository.ProductRepository  = (*mockProductRepo)(nil)
	_ repository.SupplierRepository = (*mockSupplierRepo)(nil)
	_ repository.ImageRepository    = (*mockImageRepo)(nil)
)
