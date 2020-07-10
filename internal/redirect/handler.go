package redirect

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/laytan/shortnr/internal/storage"
)

// Handler redirects request to original url or 404
func Handler(URLS storage.Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		// Check for the id in our urls map and redirect if present
		if redirectURL, ok := URLS.Get(id); ok {
			http.Redirect(w, r, redirectURL, 308)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
