package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
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

func buildMultipart(t *testing.T, productID string, imageData []byte) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	require.NoError(t, w.WriteField("product_id", productID))

	fw, err := w.CreateFormFile("image", "img.bin")
	require.NoError(t, err)
	_, err = fw.Write(imageData)
	require.NoError(t, err)

	require.NoError(t, w.Close())
	return body, w.FormDataContentType()
}

func TestImageHandler_UploadImage_Ok(t *testing.T) {
	imgID := uuid.New()
	prodID := uuid.New()

	repo := &mockImageRepo{
		createFn: func(_ context.Context, productID uuid.UUID, data []byte) (*models.Image, error) {
			assert.Equal(t, prodID, productID)
			assert.Equal(t, []byte{1, 2, 3}, data)
			return &models.Image{ID: imgID, Image: data}, nil
		},
	}
	h := NewImageHandler(repo)

	body, ct := buildMultipart(t, prodID.String(), []byte{1, 2, 3})
	req := httptest.NewRequest(http.MethodPost, "/images", body)
	req.Header.Set("Content-Type", ct)
	rec := httptest.NewRecorder()

	h.UploadImage(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp dto.ImageResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, imgID.String(), resp.ID)
	assert.Equal(t, prodID.String(), resp.ProductID)
}

func TestImageHandler_UploadImage_BadProductID(t *testing.T) {
	h := NewImageHandler(&mockImageRepo{})

	body, ct := buildMultipart(t, "not-a-uuid", []byte{1})
	req := httptest.NewRequest(http.MethodPost, "/images", body)
	req.Header.Set("Content-Type", ct)
	rec := httptest.NewRecorder()

	h.UploadImage(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestImageHandler_GetImage_Ok(t *testing.T) {
	id := uuid.New()
	data := []byte{42, 42}
	repo := &mockImageRepo{
		getByIDFn: func(_ context.Context, gotID uuid.UUID) (*models.Image, error) {
			assert.Equal(t, id, gotID)
			return &models.Image{ID: id, Image: data}, nil
		},
	}
	h := NewImageHandler(repo)

	req := newRequestWithVars(http.MethodGet, "/images/"+id.String(), "", map[string]string{"id": id.String()})
	rec := httptest.NewRecorder()
	h.GetImage(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/octet-stream", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Header().Get("Content-Disposition"), "attachment")
	assert.Equal(t, data, rec.Body.Bytes())
}

func TestImageHandler_GetImage_NotFound(t *testing.T) {
	id := uuid.New()
	repo := &mockImageRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*models.Image, error) {
			return nil, repository.ErrNotFound
		},
	}
	h := NewImageHandler(repo)

	req := newRequestWithVars(http.MethodGet, "/images/"+id.String(), "", map[string]string{"id": id.String()})
	rec := httptest.NewRecorder()
	h.GetImage(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestImageHandler_GetProductImage_Ok(t *testing.T) {
	prodID := uuid.New()
	repo := &mockImageRepo{
		getByProductIDFn: func(_ context.Context, gotID uuid.UUID) (*models.Image, error) {
			assert.Equal(t, prodID, gotID)
			return &models.Image{ID: uuid.New(), Image: []byte{1}}, nil
		},
	}
	h := NewImageHandler(repo)

	req := newRequestWithVars(http.MethodGet, "/products/"+prodID.String()+"/image", "", map[string]string{"id": prodID.String()})
	rec := httptest.NewRecorder()
	h.GetProductImage(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/octet-stream", rec.Header().Get("Content-Type"))
}

func TestImageHandler_UpdateImage_EmptyBody(t *testing.T) {
	id := uuid.New()
	h := NewImageHandler(&mockImageRepo{})

	req := newRequestWithVars(http.MethodPut, "/images/"+id.String(), "", map[string]string{"id": id.String()})
	rec := httptest.NewRecorder()
	h.UpdateImage(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestImageHandler_DeleteImage_Ok(t *testing.T) {
	id := uuid.New()
	repo := &mockImageRepo{
		deleteFn: func(_ context.Context, gotID uuid.UUID) error {
			assert.Equal(t, id, gotID)
			return nil
		},
	}
	h := NewImageHandler(repo)

	req := newRequestWithVars(http.MethodDelete, "/images/"+id.String(), "", map[string]string{"id": id.String()})
	rec := httptest.NewRecorder()
	h.DeleteImage(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
