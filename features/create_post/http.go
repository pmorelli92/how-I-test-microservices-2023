package create_post

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type rq struct {
	UserID  string `json:"userId"`
	Content string `json:"content"`
}

type rs struct {
	ID      string `json:"id"`
	UserID  string `json:"userId"`
	Content string `json:"content"`
}

type httpHandler struct {
	handler handler
}

func (h httpHandler) createPost(w http.ResponseWriter, r *http.Request) {
	var rq rq
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
		writeErrorResponse(ctx, w, err.Error(), http.StatusBadRequest)
		return
	}

	// In real world userID will be obtained from the auth token.
	post, err := h.handler.handle(ctx, rq.UserID, rq.Content)
	if err != nil {
		writeErrorResponse(ctx, w, err.Error(), http.StatusBadRequest)
		return
	}

	writeResponse(w, rs{
		ID:      post.ID,
		UserID:  post.UserID,
		Content: post.Content,
	}, http.StatusCreated)
}

func writeResponse(w http.ResponseWriter, body any, statusCode int) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, message string, statusCode int) {
	log.Println(message)
	writeResponse(w, ErrorResponse{
		Message: message,
	}, statusCode)
}
