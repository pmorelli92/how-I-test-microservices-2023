package create_post

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(router *chi.Mux, db *pgxpool.Pool) {
	repository := repository{db: db}
	handler := handler{repository: repository}
	httpHandler := httpHandler{handler: handler}

	router.Post("/post", httpHandler.createPost)
}
