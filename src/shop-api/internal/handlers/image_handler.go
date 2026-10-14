package handlers

import (
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"shop-api/internal/dto"
	"shop-api/internal/repository"
)

type ImageHandler struct {
	repo repository.ImageRepository
}

func NewImageHandler(repo repository.ImageRepository) *ImageHandler {
	return &ImageHandler{repo: repo}
}

// UploadImage godoc
// @Summary      Загрузить изображение для товара
// @Description  Принимает multipart/form-data: product_id (UUID) и файл image.
// @Tags         images
// @Accept       multipart/form-data
// @Produce      json
// @Param        product_id formData string true "UUID товара"
// @Param        image      formData file   true "Файл изображения"
// @Success      201 {object} dto.ImageResponse "Изображение загружено"
// @Failure      400 {object} dto.ErrorResponse "Некорректный product_id или пустое изображение"
// @Failure      404 {object} dto.ErrorResponse "Товар не найден"
// @Router       /images [post]
func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}
	productID, err := uuid.Parse(r.FormValue("product_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product_id")
		return
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		writeError(w, http.StatusBadRequest, "image file is required")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil || len(data) == 0 {
		writeError(w, http.StatusBadRequest, "empty image")
		return
	}

	img, err := h.repo.Create(r.Context(), productID, data)
	if err != nil {
		handleRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, dto.ImageResponse{
		ID:        img.ID.String(),
		ProductID: productID.String(),
	})
}

// GetImage godoc
// @Summary      Получить изображение по id
// @Description  Возвращает бинарный поток изображения. Файл скачивается автоматически.
// @Tags         images
// @Produce      application/octet-stream
// @Param        id path string true "UUID изображения"
// @Success      200 {string} string "Бинарные данные изображения"
// @Failure      400 {object} dto.ErrorResponse "Некорректный UUID"
// @Failure      404 {object} dto.ErrorResponse "Изображение не найдено"
// @Router       /images/{id} [get]
func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	img, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		handleRepoError(w, err)
		return
	}
	writeImage(w, img.Image)
}

// GetProductImage godoc
// @Summary      Получить изображение товара
// @Description  Возвращает бинарный поток изображения, привязанного к товару. Файл скачивается автоматически.
// @Tags         images
// @Produce      application/octet-stream
// @Param        id path string true "UUID товара"
// @Success      200 {string} string "Бинарные данные изображения"
// @Failure      400 {object} dto.ErrorResponse "Некорректный UUID"
// @Failure      404 {object} dto.ErrorResponse "Изображение не найдено"
// @Router       /products/{id}/image [get]
func (h *ImageHandler) GetProductImage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	img, err := h.repo.GetByProductID(r.Context(), id)
	if err != nil {
		handleRepoError(w, err)
		return
	}
	writeImage(w, img.Image)
}

// UpdateImage godoc
// @Summary      Изменить изображение
// @Description  Заменяет байты изображения по его id.
// @Tags         images
// @Accept       application/octet-stream
// @Produce      json
// @Param        id    path string true "UUID изображения"
// @Param        image body   []byte true "Новые байты изображения"
// @Success      200 {object} dto.ErrorResponse "Изображение обновлено"
// @Failure      400 {object} dto.ErrorResponse "Некорректный UUID или пустое тело"
// @Failure      404 {object} dto.ErrorResponse "Изображение не найдено"
// @Router       /images/{id} [put]
func (h *ImageHandler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil || len(data) == 0 {
		writeError(w, http.StatusBadRequest, "empty image body")
		return
	}
	if err := h.repo.Update(r.Context(), id, data); err != nil {
		handleRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// DeleteImage godoc
// @Summary      Удалить изображение
// @Tags         images
// @Param        id path string true "UUID изображения"
// @Success      204 "Изображение удалено"
// @Failure      400 {object} dto.ErrorResponse "Некорректный UUID"
// @Failure      404 {object} dto.ErrorResponse "Изображение не найдено"
// @Router       /images/{id} [delete]
func (h *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
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

func writeImage(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="image.bin"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
