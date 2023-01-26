package accept_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pmorelli92/how-i-test-microservices-2023/accept_test/helpers"
	"github.com/pmorelli92/how-i-test-microservices-2023/app"
)

type testCreatePost struct {
	db          *pgxpool.Pool
	HTTPAddress string
	userID      string
	content     string
}

type postRs struct {
	ID string `json:"id"`
}

type testCreatePostResult struct {
	response   postRs
	statusCode int
}

func (t testCreatePost) userExists(ctx context.Context) error {
	_, err := t.db.Exec(ctx, `INSERT INTO users(id) VALUES($1)`, t.userID)
	return err
}

func (t testCreatePost) userDoesNotExists(ctx context.Context) error {
	return nil
}

func (t testCreatePost) userCreatesPost(ctx context.Context) (testCreatePostResult, error) {
	type postRq struct {
		UserID  string `json:"userId"`
		Content string `json:"content"`
	}

	URL := fmt.Sprintf("http://%s/post", t.HTTPAddress)
	rqBody := postRq{
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

	var postRs postRs
	if err := json.Unmarshal(rsBody, &postRs); err != nil {
		return testCreatePostResult{}, err
	}

	return testCreatePostResult{
		response:   postRs,
		statusCode: rs.StatusCode,
	}, nil
}

func (t testCreatePost) responseShouldBeStatus201(ctx context.Context, result testCreatePostResult) error {
	if result.statusCode != http.StatusCreated {
		return fmt.Errorf("expected status 201, actual: %d", result.statusCode)
	}
	return nil
}

func (t testCreatePost) responseShouldBeStatus400(ctx context.Context, result testCreatePostResult) error {
	if result.statusCode != http.StatusBadRequest {
		return fmt.Errorf("expected status 400, actual: %d", result.statusCode)
	}
	return nil
}

func (t testCreatePost) postShouldExistsOnDatabase(ctx context.Context, result testCreatePostResult) error {
	row := t.db.QueryRow(
		ctx,
		`SELECT content FROM posts WHERE id = $1`, result.response.ID)

	var content string
	if err := row.Scan(&content); err != nil {
		return err
	}

	if t.content != content {
		return fmt.Errorf("expected content %s, actual: %s", t.content, content)
	}

	return nil
}

func (t testCreatePost) postIDShouldBeEmpty(ctx context.Context, result testCreatePostResult) error {
	if result.response.ID != "" {
		return fmt.Errorf("expected post ID to be empty")
	}
	return nil
}

func TestCreatePost(t *testing.T) {
	ctx := context.TODO()
	cfg, err := app.NewConfig()
	if err != nil {
		t.Fatal("could not get config")
	}

	db, err := pgxpool.New(ctx, cfg.ConnectDatabaseDSN())
	if err != nil {
		t.Fatal("could not connect to database")
	}

	t.Run("Create post without blacklisted terms", func(t *testing.T) {
		arg := testCreatePost{
			db:          db,
			HTTPAddress: cfg.HTTPAddress,
			userID:      uuid.NewString(),
			content:     "some text without blacklist",
		}

		helpers.NewCase[testCreatePostResult](ctx, t).
			Given(arg.userExists).
			When(arg.userCreatesPost).
			Then(arg.responseShouldBeStatus201).
			Then(arg.postShouldExistsOnDatabase).
			Run()
	})

	t.Run("Create post with blacklisted terms", func(t *testing.T) {
		arg := testCreatePost{
			db:          db,
			HTTPAddress: cfg.HTTPAddress,
			userID:      uuid.NewString(),
			content:     "some text with foo",
		}

		helpers.NewCase[testCreatePostResult](ctx, t).
			Given(arg.userExists).
			When(arg.userCreatesPost).
			Then(arg.responseShouldBeStatus400).
			Then(arg.postIDShouldBeEmpty).
			Run()
	})

	t.Run("Create post for non existing user", func(t *testing.T) {
		arg := testCreatePost{
			db:          db,
			HTTPAddress: cfg.HTTPAddress,
			userID:      uuid.NewString(),
			content:     "some text for non existing user",
		}

		helpers.NewCase[testCreatePostResult](ctx, t).
			Given(arg.userDoesNotExists).
			When(arg.userCreatesPost).
			Then(arg.responseShouldBeStatus400).
			Then(arg.postIDShouldBeEmpty).
			Run()
	})
}
