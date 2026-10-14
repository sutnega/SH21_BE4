package handlers

import (
	"net/http"

	authpb "auth-service/api/proto/auth/v1"

	"shop-api/internal/authclient"
	"shop-api/internal/dto"
)

type AuthHandler struct {
	auth *authclient.Client
}

func NewAuthHandler(auth *authclient.Client) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Register godoc
// @Summary      Регистрация пользователя
// @Description  Создаёт нового пользователя через сервис авторизации и возвращает JWT-токен.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body dto.RegisterRequest true "Данные пользователя"
// @Success      201 {object} dto.TokenResponse "Пользователь создан, токен выдан"
// @Failure      400 {object} dto.ErrorResponse "Некорректное тело или ошибка валидации"
// @Failure      409 {object} dto.ErrorResponse "Пользователь с таким email уже существует"
// @Failure      500 {object} dto.ErrorResponse "Ошибка сервиса авторизации"
// @Router       /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.auth.Register(r.Context(), &authpb.RegisterRequest{
		Email:    req.Email,
		Name:     req.Name,
		Surname:  req.Surname,
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		writeGrpcError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.TokenResponse{
		Token:  resp.GetToken(),
		UserID: resp.GetUserId(),
	})
}

// Login godoc
// @Summary      Аутентификация пользователя
// @Description  Проверяет пару email/пароль через сервис авторизации и возвращает JWT-токен.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        creds body dto.LoginRequest true "Учётные данные"
// @Success      200 {object} dto.TokenResponse "Токен выдан"
// @Failure      400 {object} dto.ErrorResponse "Некорректное тело или ошибка валидации"
// @Failure      401 {object} dto.ErrorResponse "Неверные учётные данные"
// @Failure      500 {object} dto.ErrorResponse "Ошибка сервиса авторизации"
// @Router       /auth [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeGrpcError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.TokenResponse{
		Token:  resp.GetToken(),
		UserID: resp.GetUserId(),
	})
}

// Reset godoc
// @Summary      Восстановление пароля
// @Description  Запускает процедуру сброса пароля через сервис авторизации.
//
//	В учебной реализации пароль выводится в консоль auth-service.
//
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body dto.ResetRequest true "Email пользователя"
// @Success      200 {object} dto.StatusResponse "Запрос принят"
// @Failure      400 {object} dto.ErrorResponse "Некорректное тело"
// @Failure      404 {object} dto.ErrorResponse "Пользователь не найден"
// @Failure      500 {object} dto.ErrorResponse "Ошибка сервиса авторизации"
// @Router       /reset [post]
func (h *AuthHandler) Reset(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.auth.ResetPassword(r.Context(), req.Email); err != nil {
		writeGrpcError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.StatusResponse{Status: "ok"})
}
