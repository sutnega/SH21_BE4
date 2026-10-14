package mappers_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shop-api/internal/dto"
	"shop-api/internal/mappers"
	"shop-api/internal/models"
)

func TestAddressDTOToModel(t *testing.T) {
	in := dto.AddressDTO{Country: "RU", City: "MSK", Street: "Tverskaya 1"}
	got := mappers.AddressDTOToModel(in)
	assert.Equal(t, "RU", got.Country)
	assert.Equal(t, "MSK", got.City)
	assert.Equal(t, "Tverskaya 1", got.Street)
}

func TestAddressModelToDTO(t *testing.T) {
	in := models.Address{Country: "DE", City: "Munich", Street: "Hauptstr 1"}
	got := mappers.AddressModelToDTO(in)
	assert.Equal(t, "DE", got.Country)
	assert.Equal(t, "Munich", got.City)
	assert.Equal(t, "Hauptstr 1", got.Street)
}

func TestClientFullToResponse(t *testing.T) {
	id := uuid.New()
	addrID := uuid.New()
	regDate := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	in := models.ClientFull{
		Client: models.Client{
			ID:               id,
			ClientName:       "Ivan",
			ClientSurname:    "Ivanov",
			Birthday:         time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
			Gender:           models.GenderMale,
			RegistrationDate: regDate,
			AddressID:        addrID,
		},
		Address: models.Address{Country: "RU", City: "MSK", Street: "Tverskaya 1"},
	}

	got := mappers.ClientFullToResponse(in)

	assert.Equal(t, id.String(), got.ID)
	assert.Equal(t, "Ivan", got.ClientName)
	assert.Equal(t, "Ivanov", got.ClientSurname)
	assert.Equal(t, "1990-05-15", got.Birthday)
	assert.Equal(t, "male", got.Gender)
	assert.Equal(t, regDate.Format(time.RFC3339), got.RegistrationDate)
	assert.Equal(t, "RU", got.Address.Country)
}

func TestProductFullToResponse_WithAndWithoutImage(t *testing.T) {
	id := uuid.New()
	supplierID := uuid.New()
	imageID := uuid.New()
	lastUpdate := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	t.Run("with image", func(t *testing.T) {
		in := models.ProductFull{
			Product: models.Product{
				ID:             id,
				Name:           "Bosch",
				Price:          100.5,
				AvailableStock: 5,
				LastUpdateDate: lastUpdate,
				SupplierID:     supplierID,
				ImageID:        &imageID,
			},
			CategoryName: "Fridges",
		}
		got := mappers.ProductFullToResponse(in)

		require.NotNil(t, got.ImageID)
		assert.Equal(t, imageID.String(), *got.ImageID)
		assert.Equal(t, "Fridges", got.Category)
		assert.Equal(t, id.String(), got.ID)
	})

	t.Run("without image", func(t *testing.T) {
		in := models.ProductFull{
			Product: models.Product{
				ID:             id,
				Name:           "Bosch",
				Price:          100.5,
				AvailableStock: 5,
				LastUpdateDate: lastUpdate,
				SupplierID:     supplierID,
				ImageID:        nil,
			},
			CategoryName: "Fridges",
		}
		got := mappers.ProductFullToResponse(in)
		assert.Nil(t, got.ImageID)
	})
}

func TestCreateProductRequestToModel(t *testing.T) {
	t.Run("valid uuid", func(t *testing.T) {
		sid := uuid.New()
		req := dto.CreateProductRequest{
			Name:           "Bosch",
			Price:          99.9,
			AvailableStock: 3,
			SupplierID:     sid.String(),
		}
		got, err := mappers.CreateProductRequestToModel(req)
		require.NoError(t, err)
		assert.Equal(t, sid, got.SupplierID)
		assert.Equal(t, "Bosch", got.Name)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		req := dto.CreateProductRequest{SupplierID: "not-a-uuid"}
		_, err := mappers.CreateProductRequestToModel(req)
		assert.Error(t, err)
	})
}

func TestSupplierFullToResponse(t *testing.T) {
	id := uuid.New()
	in := models.SupplierFull{
		Supplier: models.Supplier{
			ID:          id,
			Name:        "Bosch",
			PhoneNumber: "+49-000",
		},
		Address: models.Address{Country: "DE", City: "Munich", Street: "Hauptstr 1"},
	}
	got := mappers.SupplierFullToResponse(in)
	assert.Equal(t, id.String(), got.ID)
	assert.Equal(t, "Bosch", got.Name)
	assert.Equal(t, "DE", got.Address.Country)
	assert.Equal(t, "+49-000", got.PhoneNumber)
}
