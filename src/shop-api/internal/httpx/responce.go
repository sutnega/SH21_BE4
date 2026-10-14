package httpx

import (
	"encoding/json"
	"net/http"
)

// ErrorBody — универсальное тело ошибки. Дублирует dto.ErrorResponse,
// чтобы httpx не зависел от dto (иначе получится цикл).
type ErrorBody struct {
	Error string `json:"error"`
}

// WriteJSON пишет JSON-ответ с заданным статусом.
func WriteJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

// WriteError пишет JSON-ошибку вида {"error": "..."}.
func WriteError(w http.ResponseWriter, code int, msg string) {
	WriteJSON(w, code, ErrorBody{Error: msg})
}
