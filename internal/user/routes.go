package user

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/laytan/shortnr/internal/user/handlers"
	"github.com/laytan/shortnr/internal/user/storage"
)

// SetRoutes adds the routes needed for users service
func SetRoutes(r *mux.Router, storage storage.Storage) {
	r.HandleFunc("/signup", handlers.Signup(storage)).Methods(http.MethodPost)
}
