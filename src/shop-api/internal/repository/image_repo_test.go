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
)

func TestImageRepo_Create(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewImageRepository(db)

	imgID := uuid.New()
	prodID := uuid.New()
	data := []byte{1, 2, 3}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO images`).
		WithArgs(data).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(imgID))
	mock.ExpectExec(`UPDATE products SET image_id`).
		WithArgs(imgID, prodID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	got, err := repo.Create(context.Background(), prodID, data)
	require.NoError(t, err)
	assert.Equal(t, imgID, got.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestImageRepo_Create_ProductNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewImageRepository(db)

	imgID := uuid.New()
	prodID := uuid.New()
	data := []byte{1}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO images`).
		WithArgs(data).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(imgID))
	mock.ExpectExec(`UPDATE products SET image_id`).
		WithArgs(imgID, prodID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	_, err := repo.Create(context.Background(), prodID, data)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestImageRepo_GetByID_Ok(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewImageRepository(db)

	id := uuid.New()
	data := []byte{9, 8, 7}
	mock.ExpectQuery(`SELECT id, image FROM images`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "image"}).AddRow(id, data))

	got, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, data, got.Image)
}

func TestImageRepo_GetByID_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewImageRepository(db)

	id := uuid.New()
	mock.ExpectQuery(`SELECT id, image FROM images`).
		WithArgs(id).WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByID(context.Background(), id)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestImageRepo_GetByProductID(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewImageRepository(db)

	prodID := uuid.New()
	imgID := uuid.New()
	data := []byte{5}

	mock.ExpectQuery(`SELECT i.id, i.image`).
		WithArgs(prodID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "image"}).AddRow(imgID, data))

	got, err := repo.GetByProductID(context.Background(), prodID)
	require.NoError(t, err)
	assert.Equal(t, imgID, got.ID)
}

func TestImageRepo_Update(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewImageRepository(db)

	id := uuid.New()
	data := []byte{1, 1, 1}
	mock.ExpectExec(`UPDATE images`).
		WithArgs(data, id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, repo.Update(context.Background(), id, data))
}

func TestImageRepo_Update_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewImageRepository(db)

	id := uuid.New()
	mock.ExpectExec(`UPDATE images`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Update(context.Background(), id, []byte{1})
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestImageRepo_Delete(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewImageRepository(db)

	id := uuid.New()
	mock.ExpectExec(`DELETE FROM images`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, repo.Delete(context.Background(), id))
}
