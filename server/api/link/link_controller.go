package link

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/laytan/shortnr/pkg/ratelimit"

	"github.com/gorilla/mux"
	"github.com/laytan/shortnr/pkg/responder"
)

// SetRoutes adds the routes needed for shortn service
func SetRoutes(r *mux.Router, store Storage) {
	linkR := r.PathPrefix("").Subrouter()
	linkR.Use(ratelimit.Middleware(ratelimit.GeneralRateLimit))

	linkR.HandleFunc("/{id}", Redirect(store)).Methods(http.MethodGet)
	linkR.HandleFunc("/shortn", create(store)).Methods(http.MethodPost)
}

// Delegates to create service but handles parsing body and sending responses
func create(store Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body Link

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			responder.Err{
				Code: http.StatusBadRequest,
				Err:  errors.New("unparsable JSON body"),
			}.Send(w)
			return
		}

		id, cErr := Create(body, store)
		if cErr != nil {
			responder.CastAndSend(cErr, w)
			return
		}

		type createResponse struct {
			ID          string `json:"id"`
			OriginalURL string `json:"originalUrl"`
		}

		responder.Res{
			Message: "succesfully created link",
			Data:    createResponse{ID: id, OriginalURL: body.URL},
		}.Send(w)
	})
}

// Redirect redirects request to original url or 404
func Redirect(URLS Storage) http.HandlerFunc {
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
