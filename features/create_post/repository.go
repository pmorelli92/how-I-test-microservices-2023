package create_post

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func (r repository) userExists(ctx context.Context, userID string) (bool, error) {
	row := r.db.QueryRow(ctx, `SELECT id FROM users WHERE id = $1`, userID)
	var id string
	err := row.Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("error fetching user: %w", err)
	}

	return true, nil
}

func (r repository) storePost(ctx context.Context, post Post) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO posts(id, user_id, content, created_at)
		VALUES($1, $2, $3, $4)`, post.ID, post.UserID, post.Content, post.CreatedAt)
	if err != nil {
		return fmt.Errorf("error storing post: %w", err)
	}
	return nil
}
