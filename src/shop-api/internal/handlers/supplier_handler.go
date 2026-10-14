package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"shop-api/internal/dto"
	"shop-api/internal/mappers"
	"shop-api/internal/models"
	"shop-api/internal/repository"
)

type SupplierHandler struct {
	repo repository.SupplierRepository
}

func NewSupplierHandler(repo repository.SupplierRepository) *SupplierHandler {
	return &SupplierHandler{repo: repo}
}

// CreateSupplier godoc
// @Summary      Добавить поставщика
// @Description  Создаёт поставщика вместе с адресом.
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Param        supplier body dto.CreateSupplierRequest true "Поставщик"
// @Success      201 {object} dto.SupplierResponse "Поставщик создан"
// @Failure      400 {object} dto.ErrorResponse    "Некорректное тело запроса или ошибка валидации"
// @Failure      500 {object} dto.ErrorResponse    "Внутренняя ошибка сервера"
// @Router       /suppliers [post]
func (h *SupplierHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSupplierRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	model := &models.Supplier{Name: req.Name, PhoneNumber: req.PhoneNumber}
	addr := mappers.AddressDTOToModel(req.Address)

	created, a, err := h.repo.Create(r.Context(), model, &addr)
	if err != nil {
		handleRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, mappers.SupplierToResponse(*created, *a))
}

// ListSuppliers godoc
// @Summary      Получить всех поставщиков
// @Tags         suppliers
// @Produce      json
// @Success      200 {array}  dto.SupplierResponse "Список поставщиков (пустой массив, если ничего нет)"
// @Failure      500 {object} dto.ErrorResponse    "Внутренняя ошибка сервера"
// @Router       /suppliers [get]
func (h *SupplierHandler) ListSuppliers(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context())
	if err != nil {
		handleRepoError(w, err)
		return
	}
	resp := make([]dto.SupplierResponse, 0, len(items))
	for _, s := range items {
		resp = append(resp, mappers.SupplierFullToResponse(s))
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetSupplier godoc
// @Summary      Получить поставщика по id
// @Tags         suppliers
// @Produce      json
// @Param        id path string true "UUID поставщика"
// @Success      200 {object} dto.SupplierResponse "Поставщик найден"
// @Failure      400 {object} dto.ErrorResponse    "Некорректный UUID"
// @Failure      404 {object} dto.ErrorResponse    "Поставщик не найден"
// @Router       /suppliers/{id} [get]
func (h *SupplierHandler) GetSupplier(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	s, a, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		handleRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, mappers.SupplierToResponse(*s, *a))
}

// DeleteSupplier godoc
// @Summary      Удалить поставщика
// @Description  Удаляет поставщика и его адрес. Если на поставщика ссылаются товары — операция завершится ошибкой 500.
// @Tags         suppliers
// @Param        id path string true "UUID поставщика"
// @Success      204 "Поставщик удалён"
// @Failure      400 {object} dto.ErrorResponse "Некорректный UUID"
// @Failure      404 {object} dto.ErrorResponse "Поставщик не найден"
// @Router       /suppliers/{id} [delete]
func (h *SupplierHandler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.repo.Delete(r.Context(), id); err != nil {
		handleRepoError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdateSupplierAddress godoc
// @Summary      Изменить адрес поставщика
// @Description  Заменяет адрес поставщика на новый.
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Param        id      path string                   true "UUID поставщика"
// @Param        address body dto.UpdateAddressRequest true "Новый адрес"
// @Success      200 {object} dto.ErrorResponse "Адрес обновлён"
// @Failure      400 {object} dto.ErrorResponse "Некорректный UUID или тело запроса"
// @Failure      404 {object} dto.ErrorResponse "Поставщик не найден"
// @Router       /suppliers/{id}/address [patch]
func (h *SupplierHandler) UpdateSupplierAddress(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req dto.UpdateAddressRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	addr := mappers.AddressDTOToModel(dto.AddressDTO(req))
	if err := h.repo.UpdateAddress(r.Context(), id, &addr); err != nil {
		handleRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}
