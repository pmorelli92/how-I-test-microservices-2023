package accept_test

import (
	"os"
	"testing"
	"time"

	"github.com/pmorelli92/how-i-test-microservices-2023/app"
)

func TestMain(m *testing.M) {
	go app.Run()
	time.Sleep(500 * time.Millisecond)
	os.Exit(m.Run())
}
