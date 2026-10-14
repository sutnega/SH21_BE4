package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shop-api/internal/dto"
	"shop-api/internal/models"
	"shop-api/internal/repository"
)

func newRequestWithVars(method, path, body string, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	return mux.SetURLVars(req, vars)
}

func TestClientHandler_CreateClient_Ok(t *testing.T) {
	id := uuid.New()
	addrID := uuid.New()
	now := time.Now()

	repo := &mockClientRepo{
		createFn: func(_ context.Context, c *models.Client, a *models.Address) (*models.Client, *models.Address, error) {
			assert.Equal(t, "Ivan", c.ClientName)
			assert.Equal(t, "RU", a.Country)
			c.ID = id
			c.RegistrationDate = now
			a.ID = addrID
			return c, a, nil
		},
	}
	h := NewClientHandler(repo)

	body := `{
		"client_name":"Ivan",
		"client_surname":"Ivanov",
		"birthday":"1990-05-15",
		"gender":"male",
		"address":{"country":"RU","city":"MSK","street":"Tverskaya 1"}
	}`
	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodPost, "/clients", body, nil)

	h.CreateClient(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp dto.ClientResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, id.String(), resp.ID)
	assert.Equal(t, "Ivan", resp.ClientName)
	assert.Equal(t, "RU", resp.Address.Country)
}

func TestClientHandler_CreateClient_BadJSON(t *testing.T) {
	h := NewClientHandler(&mockClientRepo{})

	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodPost, "/clients", "{not json", nil)
	h.CreateClient(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid json")
}

func TestClientHandler_CreateClient_ValidationFails(t *testing.T) {
	h := NewClientHandler(&mockClientRepo{})

	body := `{"client_name":"","client_surname":"Ivanov","birthday":"1990-05-15","gender":"male","address":{"country":"RU","city":"MSK","street":"Tverskaya 1"}}`
	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodPost, "/clients", body, nil)
	h.CreateClient(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestClientHandler_CreateClient_BadBirthdayFormat(t *testing.T) {
	h := NewClientHandler(&mockClientRepo{})

	body := `{"client_name":"Ivan","client_surname":"Ivanov","birthday":"15/05/1990","gender":"male","address":{"country":"RU","city":"MSK","street":"Tverskaya 1"}}`
	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodPost, "/clients", body, nil)
	h.CreateClient(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Birthday")
}

func TestClientHandler_ListClients_Ok(t *testing.T) {
	id := uuid.New()
	addrID := uuid.New()

	repo := &mockClientRepo{
		listFn: func(_ context.Context, name, surname string, limit, offset int) ([]models.ClientFull, error) {
			assert.Equal(t, "Ivan", name)
			assert.Equal(t, "Ivanov", surname)
			assert.Equal(t, 5, limit)
			assert.Equal(t, 10, offset)
			return []models.ClientFull{
				{
					Client: models.Client{
						ID: id, ClientName: "Ivan", ClientSurname: "Ivanov",
						Birthday: time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
						Gender:   models.GenderMale,
					},
					Address: models.Address{ID: addrID, Country: "RU", City: "MSK", Street: "Tverskaya 1"},
				},
			}, nil
		},
	}
	h := NewClientHandler(repo)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/clients?name=Ivan&surname=Ivanov&limit=5&offset=10", nil)
	h.ListClients(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp []dto.ClientResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp, 1)
	assert.Equal(t, "Ivan", resp[0].ClientName)
}

func TestClientHandler_ListClients_EmptyReturnsEmptyArray(t *testing.T) {
	repo := &mockClientRepo{
		listFn: func(_ context.Context, _, _ string, _, _ int) ([]models.ClientFull, error) {
			return []models.ClientFull{}, nil
		},
	}
	h := NewClientHandler(repo)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/clients", nil)
	h.ListClients(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "[]\n", rec.Body.String())
}

func TestClientHandler_DeleteClient_Ok(t *testing.T) {
	id := uuid.New()
	repo := &mockClientRepo{
		deleteFn: func(_ context.Context, gotID uuid.UUID) error {
			assert.Equal(t, id, gotID)
			return nil
		},
	}
	h := NewClientHandler(repo)

	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodDelete, "/clients/"+id.String(), "", map[string]string{"id": id.String()})
	h.DeleteClient(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestClientHandler_DeleteClient_BadUUID(t *testing.T) {
	h := NewClientHandler(&mockClientRepo{})

	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodDelete, "/clients/abc", "", map[string]string{"id": "abc"})
	h.DeleteClient(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestClientHandler_DeleteClient_NotFound(t *testing.T) {
	id := uuid.New()
	repo := &mockClientRepo{
		deleteFn: func(_ context.Context, _ uuid.UUID) error { return repository.ErrNotFound },
	}
	h := NewClientHandler(repo)

	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodDelete, "/clients/"+id.String(), "", map[string]string{"id": id.String()})
	h.DeleteClient(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestClientHandler_DeleteClient_InternalError(t *testing.T) {
	id := uuid.New()
	repo := &mockClientRepo{
		deleteFn: func(_ context.Context, _ uuid.UUID) error { return errors.New("boom") },
	}
	h := NewClientHandler(repo)

	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodDelete, "/clients/"+id.String(), "", map[string]string{"id": id.String()})
	h.DeleteClient(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestClientHandler_UpdateAddress_Ok(t *testing.T) {
	id := uuid.New()
	repo := &mockClientRepo{
		updateAddressFn: func(_ context.Context, gotID uuid.UUID, a *models.Address) error {
			assert.Equal(t, id, gotID)
			assert.Equal(t, "SPb", a.City)
			return nil
		},
	}
	h := NewClientHandler(repo)

	body := `{"country":"RU","city":"SPb","street":"Nevsky 1"}`
	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodPatch, "/clients/"+id.String()+"/address", body, map[string]string{"id": id.String()})
	h.UpdateClientAddress(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestClientHandler_UpdateAddress_NotFound(t *testing.T) {
	id := uuid.New()
	repo := &mockClientRepo{
		updateAddressFn: func(_ context.Context, _ uuid.UUID, _ *models.Address) error {
			return repository.ErrNotFound
		},
	}
	h := NewClientHandler(repo)

	body := `{"country":"RU","city":"SPb","street":"Nevsky 1"}`
	rec := httptest.NewRecorder()
	req := newRequestWithVars(http.MethodPatch, "/clients/"+id.String()+"/address", body, map[string]string{"id": id.String()})
	h.UpdateClientAddress(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
