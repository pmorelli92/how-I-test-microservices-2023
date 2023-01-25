package create_post

import (
	"context"
	"fmt"
)

type handler struct {
	repository repository
}

func (h handler) handle(ctx context.Context, userID, content string) (Post, error) {
	userExists, err := h.repository.userExists(ctx, userID)
	if err != nil {
		return Post{}, nil
	}

	if !userExists {
		return Post{}, fmt.Errorf("user does not exists")
	}

	post, err := newPost(userID, content)
	if err != nil {
		return Post{}, err
	}

	if err = h.repository.storePost(ctx, post); err != nil {
		return Post{}, err
	}

	return post, nil
}
