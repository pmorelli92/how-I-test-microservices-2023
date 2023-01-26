//go:build accept

package accept_test

import (
	"os"
	"testing"
	"time"

	"github.com/pmorelli92/how-i-test-microservices-2023/app"
)

func TestMain(m *testing.M) {
	go app.Run()
	readinessProbe()
	os.Exit(m.Run())
}

// readinessProbe is used to add required ready checks for service and dependencies
func readinessProbe() {
	time.Sleep(500 * time.Millisecond)
}
