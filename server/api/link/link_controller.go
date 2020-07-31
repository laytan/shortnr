package link

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/laytan/shortnr/api/user"
	"github.com/laytan/shortnr/pkg/ratelimit"

	"github.com/gorilla/mux"
	"github.com/laytan/shortnr/pkg/responder"
)

// SetRoutes adds the routes needed for shortn service
func SetRoutes(r *mux.Router, store Storage) {
	linkR := r.PathPrefix("").Subrouter()
	linkR.Use(ratelimit.Middleware(ratelimit.GeneralRateLimit))

	linkR.HandleFunc("/{id}", one(store)).Methods(http.MethodGet)

	withAuthR := linkR.PathPrefix("").Subrouter()
	withAuthR.Use(user.JwtAuthorization)
	withAuthR.HandleFunc("", create(store)).Methods(http.MethodPost)
	withAuthR.HandleFunc("", all(store)).Methods(http.MethodGet)
}

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

		user, err := user.GetUser(r)
		if err != nil {
			responder.CastAndSend(err, w)
			return
		}
		body.UserID = user.ID

		link, cErr := Create(body, store)
		if cErr != nil {
			responder.CastAndSend(cErr, w)
			return
		}

		responder.Res{
			Message: "succesfully created link",
			Data:    link,
		}.Send(w)
	})
}

func one(store Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		if link, ok := store.Get(id); ok {
			responder.Res{
				Data: link,
			}.Send(w)
			return
		}

		responder.Err{
			Code: http.StatusNotFound,
			Err:  errors.New("link not found"),
		}.Send(w)
		return
	})
}

func all(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := user.GetUser(r)
		if err != nil {
			responder.Err{
				Code: http.StatusUnauthorized,
				Err:  errors.New("could not retrieve user"),
			}.Send(w)
			return
		}

		links := store.GetLinksFromUser(user.ID)

		responder.Res{
			Data: links,
		}.Send(w)
		return
	}
}
