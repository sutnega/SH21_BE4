package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shop-api/internal/models"
)

func TestSupplierRepo_Create(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewSupplierRepository(db)

	addrID := uuid.New()
	supID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO addresses`).
		WithArgs("DE", "Munich", "Hauptstr 1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(addrID))
	mock.ExpectQuery(`INSERT INTO suppliers`).
		WithArgs("Bosch", addrID, "+49-000").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(supID))
	mock.ExpectCommit()

	s, a, err := repo.Create(context.Background(),
		&models.Supplier{Name: "Bosch", PhoneNumber: "+49-000"},
		&models.Address{Country: "DE", City: "Munich", Street: "Hauptstr 1"},
	)
	require.NoError(t, err)
	assert.Equal(t, supID, s.ID)
	assert.Equal(t, addrID, a.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSupplierRepo_GetByID_Ok(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewSupplierRepository(db)

	id := uuid.New()
	addrID := uuid.New()

	mock.ExpectQuery(`SELECT s.id`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{
			"s_id", "name", "address_id", "phone_number",
			"a_id", "country", "city", "street",
		}).AddRow(id, "Bosch", addrID, "+49-000", addrID, "DE", "Munich", "Hauptstr 1"))

	s, a, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, "Bosch", s.Name)
	assert.Equal(t, "Munich", a.City)
}

func TestSupplierRepo_GetByID_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewSupplierRepository(db)

	id := uuid.New()
	mock.ExpectQuery(`SELECT s.id`).WithArgs(id).WillReturnError(sql.ErrNoRows)

	_, _, err := repo.GetByID(context.Background(), id)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestSupplierRepo_List(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewSupplierRepository(db)

	id := uuid.New()
	addrID := uuid.New()

	mock.ExpectQuery(`SELECT s.id`).
		WillReturnRows(sqlmock.NewRows([]string{
			"s_id", "name", "address_id", "phone_number",
			"a_id", "country", "city", "street",
		}).AddRow(id, "Bosch", addrID, "+49-000", addrID, "DE", "Munich", "Hauptstr 1"))

	got, err := repo.List(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Bosch", got[0].Name)
}

func TestSupplierRepo_UpdateAddress(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewSupplierRepository(db)

	id := uuid.New()
	mock.ExpectExec(`UPDATE addresses`).
		WithArgs("DE", "Berlin", "Alexanderplatz 1", id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, repo.UpdateAddress(context.Background(), id,
		&models.Address{Country: "DE", City: "Berlin", Street: "Alexanderplatz 1"}))
}

func TestSupplierRepo_Delete_Ok(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewSupplierRepository(db)

	id := uuid.New()
	addrID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT address_id FROM suppliers`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"address_id"}).AddRow(addrID))
	mock.ExpectExec(`DELETE FROM suppliers`).WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`DELETE FROM addresses`).WithArgs(addrID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	require.NoError(t, repo.Delete(context.Background(), id))
	require.NoError(t, mock.ExpectationsWereMet())
}
