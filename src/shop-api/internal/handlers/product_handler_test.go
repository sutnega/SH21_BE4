package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shop-api/internal/dto"
	"shop-api/internal/models"
	"shop-api/internal/repository"
)

func TestProductHandler_CreateProduct_Ok(t *testing.T) {
	id := uuid.New()
	supplierID := uuid.New()
	now := time.Now()

	repo := &mockProductRepo{
		createFn: func(_ context.Context, p *models.Product, category string) (*models.Product, error) {
			assert.Equal(t, "Bosch", p.Name)
			assert.Equal(t, "Fridges", category)
			p.ID = id
			p.LastUpdateDate = now
			return p, nil
		},
	}
	h := NewProductHandler(repo)

	body := `{
		"name":"Bosch",
		"category":"Fridges",
		"price":100.0,
		"available_stock":5,
		"supplier_id":"` + supplierID.String() + `"
	}`
	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodPost, "/products", body, nil)
	h.CreateProduct(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp dto.ProductResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, id.String(), resp.ID)
	assert.Equal(t, "Fridges", resp.Category)
}

func TestProductHandler_CreateProduct_BadSupplierID(t *testing.T) {
	h := NewProductHandler(&mockProductRepo{})

	// supplier_id проходит валидацию uuid, но кладём невалидный, чтобы проверить 400
	body := `{"name":"Bosch","category":"Fridges","price":1,"available_stock":1,"supplier_id":"abc"}`
	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodPost, "/products", body, nil)
	h.CreateProduct(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProductHandler_ListProducts_EmptyIsArray(t *testing.T) {
	repo := &mockProductRepo{
		listAvailableFn: func(_ context.Context, _, _ int) ([]models.ProductFull, error) {
			return []models.ProductFull{}, nil
		},
	}
	h := NewProductHandler(repo)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	h.ListProducts(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "[]\n", rec.Body.String())
}

func TestProductHandler_GetProduct_Ok(t *testing.T) {
	id := uuid.New()
	supplierID := uuid.New()
	now := time.Now()

	repo := &mockProductRepo{
		getByIDFn: func(_ context.Context, gotID uuid.UUID) (*models.ProductFull, error) {
			assert.Equal(t, id, gotID)
			return &models.ProductFull{
				Product: models.Product{
					ID: id, Name: "Bosch", Price: 100.0, AvailableStock: 5,
					LastUpdateDate: now, SupplierID: supplierID,
				},
				CategoryName: "Fridges",
			}, nil
		},
	}
	h := NewProductHandler(repo)

	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodGet, "/products/"+id.String(), "", map[string]string{"id": id.String()})
	h.GetProduct(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestProductHandler_GetProduct_NotFound(t *testing.T) {
	id := uuid.New()
	repo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*models.ProductFull, error) {
			return nil, repository.ErrNotFound
		},
	}
	h := NewProductHandler(repo)

	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodGet, "/products/"+id.String(), "", map[string]string{"id": id.String()})
	h.GetProduct(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestProductHandler_DecreaseStock_Ok(t *testing.T) {
	id := uuid.New()
	repo := &mockProductRepo{
		decreaseStockFn: func(_ context.Context, gotID uuid.UUID, amount int) error {
			assert.Equal(t, id, gotID)
			assert.Equal(t, 3, amount)
			return nil
		},
	}
	h := NewProductHandler(repo)

	body := `{"amount":3}`
	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodPatch, "/products/"+id.String()+"/stock", body, map[string]string{"id": id.String()})
	h.DecreaseStock(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestProductHandler_DecreaseStock_NotEnough(t *testing.T) {
	id := uuid.New()
	repo := &mockProductRepo{
		decreaseStockFn: func(_ context.Context, _ uuid.UUID, _ int) error {
			return errNotEnoughStock
		},
	}
	h := NewProductHandler(repo)

	body := `{"amount":100}`
	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodPatch, "/products/"+id.String()+"/stock", body, map[string]string{"id": id.String()})
	h.DecreaseStock(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// вспомогательная ошибка, чтобы не тянуть текстовку из репозитория в этот тест
var errNotEnoughStock = &simpleErr{"not enough stock"}

type simpleErr struct{ msg string }

func (e *simpleErr) Error() string { return e.msg }

func TestProductHandler_DeleteProduct_Ok(t *testing.T) {
	id := uuid.New()
	repo := &mockProductRepo{
		deleteFn: func(_ context.Context, gotID uuid.UUID) error {
			assert.Equal(t, id, gotID)
			return nil
		},
	}
	h := NewProductHandler(repo)

	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodDelete, "/products/"+id.String(), "", map[string]string{"id": id.String()})
	h.DeleteProduct(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
