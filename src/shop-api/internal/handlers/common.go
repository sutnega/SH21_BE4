package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"shop-api/internal/httpx"
	"shop-api/internal/repository"
)

var validate = validator.New()

// writeJSON / writeError — тонкие обёртки над httpx,
// чтобы не менять все вызовы в существующих хендлерах.
func writeJSON(w http.ResponseWriter, code int, body any) {
	httpx.WriteJSON(w, code, body)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	httpx.WriteError(w, code, msg)
}

func decodeAndValidate(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return errors.New("invalid json body")
	}
	return validate.Struct(dst)
}

func handleRepoError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "resource not found")
		return
	}
	writeError(w, http.StatusInternalServerError, err.Error())
}

func parseIntQuery(r *http.Request, key string) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return 0
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return 0
	}
	return n
}
