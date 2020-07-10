package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/laytan/shortnr/internal/storage"
	"github.com/laytan/shortnr/pkg/link"
	"github.com/rs/xid"
)

type createBody struct {
	URL string
}

// Create creates new Id to Url mapping
func Create(storage storage.Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body createBody

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "Invalid request, make sure to include a url"}`))
			return
		}

		if !link.IsLink(body.URL) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "Invalid url, make sure to include https:// or http://"}`))
			return
		}

		// Generate a unique ID
		id := xid.New().String()

		// Store in our storage
		storage.Set(id, body.URL)

		w.Write([]byte(fmt.Sprintf(`{"message": "URL created", "id": %q, "originalURL": %q}`, id, body.URL)))
	})
}
