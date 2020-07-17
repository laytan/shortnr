package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/laytan/shortnr/pkg/responder"
)

// SetRoutes adds the user related routes to the router
func SetRoutes(r *mux.Router, store Storage) {
	r.HandleFunc("/signup", signup(store)).Methods(http.MethodPost)
	r.HandleFunc("/login", login(store)).Methods(http.MethodPost)

	authR := r.PathPrefix("").Subrouter()
	authR.Use(JwtAuthorization)
	authR.HandleFunc("/me", me()).Methods(http.MethodGet)
}

func signup(store Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse body
		var credentials Credentials
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			responder.Err{
				Code: http.StatusBadRequest,
				Err:  errors.New("unable to parse JSON body"),
			}.Send(w)
			return
		}

		err = Signup(credentials, store)
		if err != nil {
			responder.CastAndSend(err, w)
			return
		}

		responder.Res{
			Message: "signed up",
		}.Send(w)
	})
}

type loginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func login(store Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse body
		var credentials Credentials
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			responder.Err{
				Code: http.StatusBadRequest,
				Err:  errors.New("unable to parse JSON body"),
			}.Send(w)
			return
		}

		// Dispatch to user service
		user, token, err := Login(credentials, store)
		if err != nil {
			responder.CastAndSend(err, w)
			return
		}

		responder.Res{
			Message: "logged in",
			Data: loginResponse{
				Token: token,
				User:  user,
			},
		}.Send(w)
	})
}

// me sends the user making the request
func me() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := GetUser(r)
		if err != nil {
			responder.CastAndSend(err, w)
			return
		}

		responder.Res{
			Data: user,
		}.Send(w)
	})
}
