package create_post

import (
	"encoding/json"
	"net/http"

	"github.com/pmorelli92/how-i-test-microservices-2023/rest"
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
		rest.WriteErrorResponse(ctx, w, err.Error(), http.StatusBadRequest)
		return
	}

	// In real world userID will be obtained from the auth token.
	post, err := h.handler.handle(ctx, rq.UserID, rq.Content)
	if err != nil {
		rest.WriteErrorResponse(ctx, w, err.Error(), http.StatusBadRequest)
		return
	}

	rest.WriteResponse(w, rs{
		ID:      post.ID,
		UserID:  post.UserID,
		Content: post.Content,
	}, http.StatusCreated)
}
