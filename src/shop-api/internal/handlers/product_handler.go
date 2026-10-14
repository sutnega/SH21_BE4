package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"shop-api/internal/dto"
	"shop-api/internal/mappers"
	"shop-api/internal/repository"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// CreateProduct godoc
// @Summary      Добавить товар
// @Description  Создаёт товар. Если категории с таким именем нет — она создаётся.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product body dto.CreateProductRequest true "Товар"
// @Success      201 {object} dto.ProductResponse "Товар создан"
// @Failure      400 {object} dto.ErrorResponse "Некорректное тело запроса или ошибка валидации"
// @Failure      500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	model, err := mappers.CreateProductRequestToModel(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid supplier_id")
		return
	}
	created, err := h.repo.Create(r.Context(), &model, req.Category)
	if err != nil {
		handleRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, mapProduct(*created, req.Category))
}

// ListProducts godoc
// @Summary      Получить все доступные товары
// @Description  Возвращает товары с available_stock > 0. Поддерживает пагинацию.
// @Tags         products
// @Produce      json
// @Param        limit  query int false "Лимит"
// @Param        offset query int false "Смещение"
// @Success      200 {array}  dto.ProductResponse "Список товаров (пустой массив, если нет доступных)"
// @Failure      500 {object} dto.ErrorResponse   "Внутренняя ошибка сервера"
// @Router       /products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	limit := parseIntQuery(r, "limit")
	offset := parseIntQuery(r, "offset")
	items, err := h.repo.ListAvailable(r.Context(), limit, offset)
	if err != nil {
		handleRepoError(w, err)
		return
	}
	resp := make([]dto.ProductResponse, 0, len(items))
	for _, p := range items {
		resp = append(resp, mappers.ProductFullToResponse(p))
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetProduct godoc
// @Summary      Получить товар по id
// @Tags         products
// @Produce      json
// @Param        id path string true "UUID товара"
// @Success      200 {object} dto.ProductResponse "Товар найден"
// @Failure      400 {object} dto.ErrorResponse   "Некорректный UUID"
// @Failure      404 {object} dto.ErrorResponse   "Товар не найден"
// @Router       /products/{id} [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	p, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		handleRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, mappers.ProductFullToResponse(*p))
}

// DeleteProduct godoc
// @Summary      Удалить товар
// @Tags         products
// @Param        id path string true "UUID товара"
// @Success      204 "Товар удалён"
// @Failure      400 {object} dto.ErrorResponse "Некорректный UUID"
// @Failure      404 {object} dto.ErrorResponse "Товар не найден"
// @Router       /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
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

// DecreaseStock godoc
// @Summary      Уменьшить количество товара
// @Description  Уменьшает available_stock на указанное число. Возвращает 400, если товара недостаточно.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id     path string                   true "UUID товара"
// @Param        amount body dto.DecreaseStockRequest true "Сколько вычесть"
// @Success      200 {object} dto.ErrorResponse "Остаток обновлён"
// @Failure      400 {object} dto.ErrorResponse "Недостаточно товара, некорректный UUID или тело"
// @Failure      404 {object} dto.ErrorResponse "Товар не найден"
// @Router       /products/{id}/stock [patch]
func (h *ProductHandler) DecreaseStock(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req dto.DecreaseStockRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.repo.DecreaseStock(r.Context(), id, req.Amount); err != nil {
		if err.Error() == "not enough stock" {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		handleRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}
