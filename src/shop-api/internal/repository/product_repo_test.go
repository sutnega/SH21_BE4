package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shop-api/internal/models"
)

func productRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "name", "category_id", "price", "available_stock",
		"last_update_date", "supplier_id", "image_id", "category_name",
	})
}

func TestProductRepo_Create_NewCategory(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewProductRepository(db)

	catID := uuid.New()
	prodID := uuid.New()
	supplierID := uuid.New()
	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT id FROM categories`).
		WithArgs("Fridges").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery(`INSERT INTO categories`).
		WithArgs("Fridges").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(catID))
	mock.ExpectQuery(`INSERT INTO products`).
		WithArgs("Bosch", catID, 100.0, 5, supplierID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "last_update_date", "image_id"}).
			AddRow(prodID, now, nil))
	mock.ExpectCommit()

	got, err := repo.Create(context.Background(), &models.Product{
		Name: "Bosch", Price: 100.0, AvailableStock: 5, SupplierID: supplierID,
	}, "Fridges")
	require.NoError(t, err)
	assert.Equal(t, prodID, got.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_Create_ExistingCategory(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewProductRepository(db)

	catID := uuid.New()
	prodID := uuid.New()
	supplierID := uuid.New()
	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT id FROM categories`).
		WithArgs("Fridges").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(catID))
	mock.ExpectQuery(`INSERT INTO products`).
		WithArgs("Bosch", catID, 100.0, 5, supplierID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "last_update_date", "image_id"}).
			AddRow(prodID, now, nil))
	mock.ExpectCommit()

	_, err := repo.Create(context.Background(), &models.Product{
		Name: "Bosch", Price: 100.0, AvailableStock: 5, SupplierID: supplierID,
	}, "Fridges")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_GetByID_Ok(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewProductRepository(db)

	id := uuid.New()
	catID := uuid.New()
	supplierID := uuid.New()
	now := time.Now()

	mock.ExpectQuery(`SELECT p.id`).
		WithArgs(id).
		WillReturnRows(productRows().AddRow(
			id, "Bosch", catID, 100.0, 5, now, supplierID, nil, "Fridges",
		))

	got, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, "Bosch", got.Name)
	assert.Equal(t, "Fridges", got.CategoryName)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_GetByID_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewProductRepository(db)

	id := uuid.New()
	mock.ExpectQuery(`SELECT p.id`).WithArgs(id).WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByID(context.Background(), id)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestProductRepo_ListAvailable(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewProductRepository(db)

	id := uuid.New()
	catID := uuid.New()
	supplierID := uuid.New()
	now := time.Now()

	mock.ExpectQuery(`SELECT p.id`).
		WithArgs(10).
		WillReturnRows(productRows().AddRow(
			id, "Bosch", catID, 100.0, 5, now, supplierID, nil, "Fridges",
		))

	got, err := repo.ListAvailable(context.Background(), 10, 0)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_DecreaseStock_Ok(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewProductRepository(db)

	id := uuid.New()
	mock.ExpectExec(`UPDATE products`).
		WithArgs(2, id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, repo.DecreaseStock(context.Background(), id, 2))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_DecreaseStock_NotEnough(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewProductRepository(db)

	id := uuid.New()
	mock.ExpectExec(`UPDATE products`).
		WithArgs(100, id).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	err := repo.DecreaseStock(context.Background(), id, 100)
	require.Error(t, err)
	assert.Equal(t, "not enough stock", err.Error())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_DecreaseStock_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewProductRepository(db)

	id := uuid.New()
	mock.ExpectExec(`UPDATE products`).
		WithArgs(1, id).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	err := repo.DecreaseStock(context.Background(), id, 1)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestProductRepo_Delete_Ok(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewProductRepository(db)

	id := uuid.New()
	mock.ExpectExec(`DELETE FROM products`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, repo.Delete(context.Background(), id))
}

func TestProductRepo_Delete_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewProductRepository(db)

	id := uuid.New()
	mock.ExpectExec(`DELETE FROM products`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Delete(context.Background(), id)
	assert.True(t, errors.Is(err, ErrNotFound))
}
