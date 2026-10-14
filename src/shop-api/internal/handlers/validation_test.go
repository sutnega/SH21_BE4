package handlers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shop-api/internal/dto"
)

func TestValidateCreateClientRequest(t *testing.T) {
	valid := dto.CreateClientRequest{
		ClientName:    "Ivan",
		ClientSurname: "Ivanov",
		Birthday:      "1990-05-15",
		Gender:        "male",
		Address:       dto.AddressDTO{Country: "RU", City: "MSK", Street: "Tverskaya 1"},
	}

	t.Run("ok", func(t *testing.T) {
		require.NoError(t, validate.Struct(valid))
	})

	t.Run("missing name", func(t *testing.T) {
		v := valid
		v.ClientName = ""
		err := validate.Struct(v)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "clientname")
	})

	t.Run("bad birthday format", func(t *testing.T) {
		v := valid
		v.Birthday = "15.05.1990"
		assert.Error(t, validate.Struct(v))
	})

	t.Run("bad gender", func(t *testing.T) {
		v := valid
		v.Gender = "attack helicopter"
		assert.Error(t, validate.Struct(v))
	})

	t.Run("empty address", func(t *testing.T) {
		v := valid
		v.Address = dto.AddressDTO{}
		assert.Error(t, validate.Struct(v))
	})
}

func TestValidateCreateProductRequest(t *testing.T) {
	valid := dto.CreateProductRequest{
		Name:           "Bosch",
		Category:       "Fridges",
		Price:          100.0,
		AvailableStock: 5,
		SupplierID:     "550e8400-e29b-41d4-a716-446655440000",
	}

	t.Run("ok", func(t *testing.T) {
		require.NoError(t, validate.Struct(valid))
	})

	t.Run("negative price", func(t *testing.T) {
		v := valid
		v.Price = -1
		assert.Error(t, validate.Struct(v))
	})

	t.Run("negative stock", func(t *testing.T) {
		v := valid
		v.AvailableStock = -1
		assert.Error(t, validate.Struct(v))
	})

	t.Run("bad uuid", func(t *testing.T) {
		v := valid
		v.SupplierID = "abc"
		assert.Error(t, validate.Struct(v))
	})
}

func TestValidateCreateSupplierRequest(t *testing.T) {
	valid := dto.CreateSupplierRequest{
		Name:        "Bosch",
		PhoneNumber: "+49-000",
		Address:     dto.AddressDTO{Country: "DE", City: "Munich", Street: "Hauptstr 1"},
	}

	t.Run("ok", func(t *testing.T) {
		require.NoError(t, validate.Struct(valid))
	})
	t.Run("empty phone", func(t *testing.T) {
		v := valid
		v.PhoneNumber = ""
		assert.Error(t, validate.Struct(v))
	})
}

func TestValidateDecreaseStockRequest(t *testing.T) {
	require.NoError(t, validate.Struct(dto.DecreaseStockRequest{Amount: 1}))
	assert.Error(t, validate.Struct(dto.DecreaseStockRequest{Amount: 0}))
	assert.Error(t, validate.Struct(dto.DecreaseStockRequest{Amount: -5}))
}
