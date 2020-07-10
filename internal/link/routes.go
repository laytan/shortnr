package link

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/laytan/shortnr/internal/link/handlers"
	"github.com/laytan/shortnr/internal/storage"
)

// SetRoutes adds the routes needed for shortn service
func SetRoutes(r *mux.Router, storage storage.Storage) {
	r.HandleFunc("/shortn", handlers.Create(storage)).Methods(http.MethodPost)
}
