package accept_test

import (
	"net/http"
	"testing"
)

func TestCreatePost(t *testing.T) {
	URL := "http://localhost:8080/dummy"
	rq, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		t.Fatal("creating rq failed")
	}

	rs, err := http.DefaultClient.Do(rq)
	if err != nil {
		t.Fatal("could not do http call rq failed")
	}

	if rs.StatusCode != http.StatusOK {
		t.Fatal("status code not expected")
	}
}
