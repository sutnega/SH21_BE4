package mappers

import (
	"time"

	"github.com/google/uuid"

	"shop-api/internal/dto"
	"shop-api/internal/models"
)

// ---------- Address ----------

func AddressDTOToModel(a dto.AddressDTO) models.Address {
	return models.Address{Country: a.Country, City: a.City, Street: a.Street}
}

func AddressModelToDTO(a models.Address) dto.AddressDTO {
	return dto.AddressDTO{Country: a.Country, City: a.City, Street: a.Street}
}

// ---------- Client ----------

func ClientFullToResponse(c models.ClientFull) dto.ClientResponse {
	return dto.ClientResponse{
		ID:               c.Client.ID.String(),
		ClientName:       c.Client.ClientName,
		ClientSurname:    c.Client.ClientSurname,
		Birthday:         c.Client.Birthday.Format("2006-01-02"),
		Gender:           string(c.Client.Gender),
		RegistrationDate: c.Client.RegistrationDate.Format(time.RFC3339),
		Address:          AddressModelToDTO(c.Address),
	}
}

func ClientToResponse(c models.Client, a models.Address) dto.ClientResponse {
	return ClientFullToResponse(models.ClientFull{Client: c, Address: a})
}

// ---------- Product ----------

func ProductFullToResponse(p models.ProductFull) dto.ProductResponse {
	var imageID *string
	if p.ImageID != nil {
		s := p.ImageID.String()
		imageID = &s
	}
	return dto.ProductResponse{
		ID:             p.ID.String(),
		Name:           p.Name,
		Category:       p.CategoryName,
		Price:          p.Price,
		AvailableStock: p.AvailableStock,
		LastUpdateDate: p.LastUpdateDate.Format(time.RFC3339),
		SupplierID:     p.SupplierID.String(),
		ImageID:        imageID,
	}
}

func CreateProductRequestToModel(r dto.CreateProductRequest) (models.Product, error) {
	supplierID, err := uuid.Parse(r.SupplierID)
	if err != nil {
		return models.Product{}, err
	}
	return models.Product{
		Name:           r.Name,
		Price:          r.Price,
		AvailableStock: r.AvailableStock,
		SupplierID:     supplierID,
	}, nil
}

// ---------- Supplier ----------

func SupplierFullToResponse(s models.SupplierFull) dto.SupplierResponse {
	return dto.SupplierResponse{
		ID:          s.Supplier.ID.String(),
		Name:        s.Supplier.Name,
		Address:     AddressModelToDTO(s.Address),
		PhoneNumber: s.Supplier.PhoneNumber,
	}
}

func SupplierToResponse(s models.Supplier, a models.Address) dto.SupplierResponse {
	return SupplierFullToResponse(models.SupplierFull{Supplier: s, Address: a})
}
