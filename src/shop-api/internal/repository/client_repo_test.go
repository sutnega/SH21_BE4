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

func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return db, mock
}

func TestClientRepo_Create(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewClientRepository(db)

	addrID := uuid.New()
	clientID := uuid.New()
	regDate := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO addresses`).
		WithArgs("RU", "MSK", "Tverskaya 1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(addrID))
	mock.ExpectQuery(`INSERT INTO clients`).
		WithArgs("Ivan", "Ivanov", sqlmock.AnyArg(), models.GenderMale, addrID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date"}).AddRow(clientID, regDate))
	mock.ExpectCommit()

	c := &models.Client{
		ClientName:    "Ivan",
		ClientSurname: "Ivanov",
		Birthday:      time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
		Gender:        models.GenderMale,
	}
	a := &models.Address{Country: "RU", City: "MSK", Street: "Tverskaya 1"}

	gotC, gotA, err := repo.Create(context.Background(), c, a)
	require.NoError(t, err)
	assert.Equal(t, clientID, gotC.ID)
	assert.Equal(t, addrID, gotA.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestClientRepo_GetByID_Ok(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewClientRepository(db)

	id := uuid.New()
	addrID := uuid.New()
	regDate := time.Now()
	birthday := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery(`SELECT c.id`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "client_name", "client_surname", "birthday", "gender",
			"registration_date", "address_id",
			"addr_id", "country", "city", "street",
		}).AddRow(
			id, "Ivan", "Ivanov", birthday, "male",
			regDate, addrID,
			addrID, "RU", "MSK", "Tverskaya 1",
		))

	c, a, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, "Ivan", c.ClientName)
	assert.Equal(t, "MSK", a.City)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestClientRepo_GetByID_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewClientRepository(db)

	id := uuid.New()
	mock.ExpectQuery(`SELECT c.id`).WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	_, _, err := repo.GetByID(context.Background(), id)
	assert.True(t, errors.Is(err, ErrNotFound))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestClientRepo_List_FiltersAndPagination(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewClientRepository(db)

	id := uuid.New()
	addrID := uuid.New()
	regDate := time.Now()
	birthday := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery(`SELECT c.id`).
		WithArgs("Ivan", "Ivanov", 10, 20).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "client_name", "client_surname", "birthday", "gender",
			"registration_date", "address_id",
			"addr_id", "country", "city", "street",
		}).AddRow(
			id, "Ivan", "Ivanov", birthday, "male",
			regDate, addrID,
			addrID, "RU", "MSK", "Tverskaya 1",
		))

	got, err := repo.List(context.Background(), "Ivan", "Ivanov", 10, 20)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Ivan", got[0].ClientName)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestClientRepo_List_EmptyReturnsEmptySlice(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewClientRepository(db)

	mock.ExpectQuery(`SELECT c.id`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "client_name", "client_surname", "birthday", "gender",
			"registration_date", "address_id",
			"addr_id", "country", "city", "street",
		}))

	got, err := repo.List(context.Background(), "", "", 0, 0)
	require.NoError(t, err)
	assert.NotNil(t, got)
	assert.Len(t, got, 0)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestClientRepo_UpdateAddress_Ok(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewClientRepository(db)

	id := uuid.New()
	mock.ExpectExec(`UPDATE addresses`).
		WithArgs("RU", "SPb", "Nevsky 1", id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateAddress(context.Background(), id, &models.Address{
		Country: "RU", City: "SPb", Street: "Nevsky 1",
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestClientRepo_UpdateAddress_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewClientRepository(db)

	id := uuid.New()
	mock.ExpectExec(`UPDATE addresses`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdateAddress(context.Background(), id, &models.Address{})
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestClientRepo_Delete_Ok(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewClientRepository(db)

	id := uuid.New()
	addrID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT address_id FROM clients`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"address_id"}).AddRow(addrID))
	mock.ExpectExec(`DELETE FROM clients`).WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`DELETE FROM addresses`).WithArgs(addrID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestClientRepo_Delete_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewClientRepository(db)

	id := uuid.New()
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT address_id FROM clients`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	err := repo.Delete(context.Background(), id)
	assert.True(t, errors.Is(err, ErrNotFound))
}
