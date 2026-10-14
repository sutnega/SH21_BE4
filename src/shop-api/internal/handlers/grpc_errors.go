package handlers

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// writeGrpcError переводит gRPC-статус в HTTP-код и пишет ответ.
func writeGrpcError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	code := http.StatusInternalServerError
	switch st.Code() {
	case codes.InvalidArgument:
		code = http.StatusBadRequest
	case codes.NotFound:
		code = http.StatusNotFound
	case codes.AlreadyExists:
		code = http.StatusConflict
	case codes.Unauthenticated, codes.PermissionDenied:
		code = http.StatusUnauthorized
	case codes.Unavailable, codes.DeadlineExceeded:
		code = http.StatusServiceUnavailable
	case codes.Internal, codes.Unknown:
		code = http.StatusInternalServerError
	}

	msg := st.Message()
	if msg == "" {
		msg = "auth service error"
	}
	writeError(w, code, msg)
}
