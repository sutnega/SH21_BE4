package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shop-api/internal/dto"
	"shop-api/internal/models"
	"shop-api/internal/repository"
)

func TestSupplierHandler_CreateSupplier_Ok(t *testing.T) {
	id := uuid.New()
	addrID := uuid.New()

	repo := &mockSupplierRepo{
		createFn: func(_ context.Context, s *models.Supplier, a *models.Address) (*models.Supplier, *models.Address, error) {
			s.ID = id
			a.ID = addrID
			return s, a, nil
		},
	}
	h := NewSupplierHandler(repo)

	body := `{"name":"Bosch","phone_number":"+49-000","address":{"country":"DE","city":"Munich","street":"Hauptstr 1"}}`
	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodPost, "/suppliers", body, nil)
	h.CreateSupplier(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp dto.SupplierResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, id.String(), resp.ID)
	assert.Equal(t, "DE", resp.Address.Country)
}

func TestSupplierHandler_CreateSupplier_ValidationFails(t *testing.T) {
	h := NewSupplierHandler(&mockSupplierRepo{})

	body := `{"name":"","phone_number":"","address":{}}`
	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodPost, "/suppliers", body, nil)
	h.CreateSupplier(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSupplierHandler_ListSuppliers_Empty(t *testing.T) {
	repo := &mockSupplierRepo{
		listFn: func(_ context.Context) ([]models.SupplierFull, error) {
			return []models.SupplierFull{}, nil
		},
	}
	h := NewSupplierHandler(repo)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/suppliers", nil)
	h.ListSuppliers(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "[]\n", rec.Body.String())
}

func TestSupplierHandler_GetSupplier_NotFound(t *testing.T) {
	id := uuid.New()
	repo := &mockSupplierRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*models.Supplier, *models.Address, error) {
			return nil, nil, repository.ErrNotFound
		},
	}
	h := NewSupplierHandler(repo)

	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodGet, "/suppliers/"+id.String(), "", map[string]string{"id": id.String()})
	h.GetSupplier(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestSupplierHandler_DeleteSupplier_Ok(t *testing.T) {
	id := uuid.New()
	repo := &mockSupplierRepo{
		deleteFn: func(_ context.Context, gotID uuid.UUID) error {
			assert.Equal(t, id, gotID)
			return nil
		},
	}
	h := NewSupplierHandler(repo)

	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodDelete, "/suppliers/"+id.String(), "", map[string]string{"id": id.String()})
	h.DeleteSupplier(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestSupplierHandler_UpdateAddress_Ok(t *testing.T) {
	id := uuid.New()
	repo := &mockSupplierRepo{
		updateAddressFn: func(_ context.Context, gotID uuid.UUID, a *models.Address) error {
			assert.Equal(t, id, gotID)
			assert.Equal(t, "Berlin", a.City)
			return nil
		},
	}
	h := NewSupplierHandler(repo)

	body := `{"country":"DE","city":"Berlin","street":"Alexanderplatz 1"}`
	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodPatch, "/suppliers/"+id.String()+"/address", body, map[string]string{"id": id.String()})
	h.UpdateSupplierAddress(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
