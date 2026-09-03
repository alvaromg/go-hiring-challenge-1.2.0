package api

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/test/helper"
)

// TestMain boots a real Postgres instance in a Docker container via gnomock,
// seeds it with the same SQL files used in production (./sql), and wires up
// the application the same way cmd/server/main.go does. All tests in this
// package share this single container and hit the HTTP router directly.
func TestCatalog(t *testing.T) {

	router, close := helper.BuildApi()
	defer close()

	testDefaultNotFound(t, router)
}
