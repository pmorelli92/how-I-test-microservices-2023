package accept_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pmorelli92/how-i-test-microservices-2023/accept_test/helpers"
)

type testCreatePost struct {
	db      *pgxpool.Pool
	userID  string
	content string
}

type testCreatePostResult struct {
	body       []byte
	statusCode int
}

func (t testCreatePost) userExists(ctx context.Context) error {
	_, err := t.db.Exec(ctx, `INSERT INTO users(id) VALUES($1)`, t.userID)
	return err
}

// func (t testCreatePost) userDoesNotExists(ctx context.Context) error {
// 	return nil
// }

func (t testCreatePost) userCreatesPost(ctx context.Context) (testCreatePostResult, error) {
	type post struct {
		UserID  string `json:"userId"`
		Content string `json:"content"`
	}

	URL := fmt.Sprintf("http://%s/post", os.Getenv("HTTP_ADDRESS"))
	rqBody := post{
		UserID:  t.userID,
		Content: t.content,
	}
	b, err := json.Marshal(rqBody)
	if err != nil {
		return testCreatePostResult{}, err
	}

	rq, err := http.NewRequest(http.MethodPost, URL, bytes.NewReader(b))
	if err != nil {
		return testCreatePostResult{}, err
	}

	rs, err := http.DefaultClient.Do(rq)
	if err != nil {
		return testCreatePostResult{}, err
	}

	defer rs.Body.Close()
	rsBody, err := io.ReadAll(rs.Body)
	if err != nil {
		return testCreatePostResult{}, err
	}

	return testCreatePostResult{
		body:       rsBody,
		statusCode: rs.StatusCode,
	}, nil
}

func (t testCreatePost) responseShouldBeStatus201(ctx context.Context, result testCreatePostResult) error {
	if result.statusCode != http.StatusCreated {
		return fmt.Errorf("expected status 201, actual: %d", result.statusCode)
	}
	return nil
}

func (t testCreatePost) postShouldExistsOnDatabase(ctx context.Context, result testCreatePostResult) error {
	type post struct {
		ID string `json:"id"`
	}

	var p post
	if err := json.Unmarshal(result.body, &p); err != nil {
		return err
	}

	row := t.db.QueryRow(
		ctx,
		`SELECT content FROM posts WHERE id = $1`, p.ID)

	var content string
	if err := row.Scan(&content); err != nil {
		return err
	}

	if t.content != content {
		return fmt.Errorf("expected content %s, actual: %s", t.content, content)
	}

	return nil
}

func databaseConn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func TestCreatePost(t *testing.T) {
	ctx := context.TODO()
	db, err := pgxpool.New(ctx, databaseConn())
	if err != nil {
		t.Fatal("could not connect to database")
	}

	arg := testCreatePost{
		db:      db,
		userID:  uuid.NewString(),
		content: "some text without blacklist",
	}

	helpers.NewCase[testCreatePostResult](ctx, t).
		Given(arg.userExists).
		When(arg.userCreatesPost).
		Then(arg.responseShouldBeStatus201).
		Then(arg.postShouldExistsOnDatabase).
		Run()
}
