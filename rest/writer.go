package rest

import (
	"context"
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, body any, statusCode int) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func WriteErrorResponse(ctx context.Context, w http.ResponseWriter, message string, statusCode int) {
	WriteResponse(w, ErrorResponse{
		Message: message,
	}, statusCode)
}
