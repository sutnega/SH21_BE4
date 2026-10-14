package middleware

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shop-api/internal/httpx"
)

// TokenVerifier — минимальный интерфейс, чтобы middleware не зависел
// от конкретного authclient (удобно для тестов).
type TokenVerifier interface {
	VerifyToken(ctx context.Context, token string) (*VerifyResult, error)
}

// VerifyResult — то, что middleware нужно знать о токене.
type VerifyResult struct {
	Valid  bool
	UserID string
	Email  string
}

type userIDKey struct{}
type emailKey struct{}

// UserID достаёт userID из контекста запроса.
func UserID(ctx context.Context) string {
	if v, ok := ctx.Value(userIDKey{}).(string); ok {
		return v
	}
	return ""
}

// Email достаёт email из контекста.
func Email(ctx context.Context) string {
	if v, ok := ctx.Value(emailKey{}).(string); ok {
		return v
	}
	return ""
}

// Auth — middleware авторизации.
func Auth(verifier TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token, ok := extractBearer(header)
			if !ok {
				httpx.WriteError(w, http.StatusUnauthorized, "missing or malformed Authorization header")
				return
			}

			res, err := verifier.VerifyToken(r.Context(), token)
			if err != nil {
				if st, ok := status.FromError(err); ok {
					switch st.Code() {
					case codes.Unauthenticated, codes.InvalidArgument, codes.NotFound:
						httpx.WriteError(w, http.StatusUnauthorized, "invalid token")
						return
					}
				}
				httpx.WriteError(w, http.StatusInternalServerError, "auth service unavailable")
				return
			}
			if res == nil || !res.Valid {
				httpx.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey{}, res.UserID)
			ctx = context.WithValue(ctx, emailKey{}, res.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractBearer(header string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", false
	}
	return token, true
}
