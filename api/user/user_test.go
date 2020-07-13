package user

import (
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/laytan/shortnr/internal/user/routes"
	"github.com/laytan/shortnr/internal/user/storage"
	"github.com/laytan/shortnr/pkg/jsonmiddleware"
)

// SetupUserServer returns a server and a store with the user endpoints registered, make sure to call Close on returned server
func SetupUserServer() (*httptest.Server, *storage.Memory) {
	// Load .env
	envErr := godotenv.Load("../../.env")
	if envErr != nil {
		panic(envErr)
	}

	// Set up store
	store := storage.Memory{Users: make([]storage.User, 0)}

	// Set up handler
	r := mux.NewRouter()
	r.Use(jsonmiddleware.Middleware)
	routes.Set(r, &store)

	// Set up server
	return httptest.NewServer(r), &store
}
