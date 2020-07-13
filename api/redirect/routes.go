package redirect

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/laytan/shortnr/internal/storage"
)

// SetRoutes adds the routes redirect uses
func SetRoutes(r *mux.Router, storage storage.Storage) {
	r.HandleFunc("/{id}", Handler(storage)).Methods(http.MethodGet)
}
