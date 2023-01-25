package create_post

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        string
	UserID    string
	Content   string
	CreatedAt time.Time
}

func newPost(userID, content string) (Post, error) {
	if err := validateContent(content); err != nil {
		return Post{}, err
	}

	return Post{
		ID:        uuid.NewString(),
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}

var blackListRegex = regexp.MustCompile(`\b(bar|foo)\b`)

func validateContent(content string) error {
	for _, match := range blackListRegex.FindAllString(content, -1) {
		return fmt.Errorf("found blacklisted word %s", match)
	}
	return nil
}
