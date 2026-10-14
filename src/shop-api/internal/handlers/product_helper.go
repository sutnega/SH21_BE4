package handlers

import (
	"time"

	"shop-api/internal/dto"
	"shop-api/internal/models"
)

// mapProduct собирает DTO-ответ после создания товара (когда CategoryName ещё не прочитан из БД,
// но мы его знаем из запроса).
func mapProduct(p models.Product, categoryName string) dto.ProductResponse {
	var imageID *string
	if p.ImageID != nil {
		s := p.ImageID.String()
		imageID = &s
	}
	return dto.ProductResponse{
		ID:             p.ID.String(),
		Name:           p.Name,
		Category:       categoryName,
		Price:          p.Price,
		AvailableStock: p.AvailableStock,
		LastUpdateDate: p.LastUpdateDate.Format(time.RFC3339),
		SupplierID:     p.SupplierID.String(),
		ImageID:        imageID,
	}
}

func structToFull(p models.Product, categoryName string) models.ProductFull {
	return models.ProductFull{Product: p, CategoryName: categoryName}
}
