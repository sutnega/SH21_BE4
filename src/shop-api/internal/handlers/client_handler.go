package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"shop-api/internal/dto"
	"shop-api/internal/mappers"
	"shop-api/internal/models"
	"shop-api/internal/repository"
)

type ClientHandler struct {
	repo repository.ClientRepository
}

func NewClientHandler(repo repository.ClientRepository) *ClientHandler {
	return &ClientHandler{repo: repo}
}

// CreateClient godoc
// @Summary      Добавить клиента
// @Description  Создаёт нового клиента вместе с адресом.
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        client body dto.CreateClientRequest true "Данные клиента"
// @Success      201 {object} dto.ClientResponse "Клиент создан"
// @Failure      400 {object} dto.ErrorResponse "Некорректное тело запроса или ошибка валидации"
// @Failure      500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /clients [post]
func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateClientRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid birthday format, expected YYYY-MM-DD")
		return
	}

	client := &models.Client{
		ClientName:    req.ClientName,
		ClientSurname: req.ClientSurname,
		Birthday:      birthday,
		Gender:        models.Gender(req.Gender),
	}
	address := mappers.AddressDTOToModel(req.Address)

	created, addr, err := h.repo.Create(r.Context(), client, &address)
	if err != nil {
		handleRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, mappers.ClientToResponse(*created, *addr))
}

// ListClients godoc
// @Summary      Получить список клиентов
// @Description  Возвращает список клиентов. Опционально фильтруется по имени и фамилии и поддерживает пагинацию.
// @Tags         clients
// @Produce      json
// @Param        name    query string false "Имя клиента"
// @Param        surname query string false "Фамилия клиента"
// @Param        limit   query int    false "Лимит"
// @Param        offset  query int    false "Смещение"
// @Success      200 {array}  dto.ClientResponse "Список клиентов (пустой массив, если ничего не найдено)"
// @Failure      500 {object} dto.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /clients [get]
func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	surname := r.URL.Query().Get("surname")
	limit := parseIntQuery(r, "limit")
	offset := parseIntQuery(r, "offset")

	items, err := h.repo.List(r.Context(), name, surname, limit, offset)
	if err != nil {
		handleRepoError(w, err)
		return
	}
	resp := make([]dto.ClientResponse, 0, len(items))
	for _, c := range items {
		resp = append(resp, mappers.ClientFullToResponse(c))
	}
	writeJSON(w, http.StatusOK, resp)
}

// DeleteClient godoc
// @Summary      Удалить клиента
// @Description  Удаляет клиента и связанный с ним адрес.
// @Tags         clients
// @Param        id path string true "UUID клиента"
// @Success      204 "Клиент удалён"
// @Failure      400 {object} dto.ErrorResponse "Некорректный UUID"
// @Failure      404 {object} dto.ErrorResponse "Клиент не найден"
// @Router       /clients/{id} [delete]
func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
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

// UpdateClientAddress godoc
// @Summary      Изменить адрес клиента
// @Description  Заменяет адрес клиента на новый.
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        id      path string                   true "UUID клиента"
// @Param        address body dto.UpdateAddressRequest true "Новый адрес"
// @Success      200 {object} dto.ErrorResponse "Адрес обновлён"
// @Failure      400 {object} dto.ErrorResponse "Некорректный UUID или тело запроса"
// @Failure      404 {object} dto.ErrorResponse "Клиент не найден"
// @Router       /clients/{id}/address [patch]
func (h *ClientHandler) UpdateClientAddress(w http.ResponseWriter, r *http.Request) {
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
