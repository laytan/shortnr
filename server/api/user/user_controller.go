package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/laytan/shortnr/pkg/ratelimit"
	"github.com/laytan/shortnr/pkg/responder"
)

type tokenResponse struct {
	Token string `json:"token"`
}

// SetRoutes adds the user related routes to the router
func SetRoutes(r *mux.Router, store Storage) {
	// Rate limit to 1 request per 10 seconds
	userR := r.PathPrefix("").Subrouter()
	userR.Use(ratelimit.Middleware(ratelimit.AuthorizationRateLimit))

	userR.HandleFunc("/signup", signup(store)).Methods(http.MethodPost)
	userR.HandleFunc("/login", login(store)).Methods(http.MethodPost)
	userR.HandleFunc("/refresh", refresh(store)).Methods(http.MethodPost)
	userR.HandleFunc("/logout", logout).Methods(http.MethodDelete)

	authR := userR.PathPrefix("").Subrouter()
	authR.Use(JwtAuthorization)

	authR.HandleFunc("/me", me).Methods(http.MethodGet)
}

// Helper to generate refresh cookie with required settings
func refreshCookie(refreshToken string) http.Cookie {
	return http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
		Secure:   true,
	}
}

func signup(store Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse body
		var credentials Credentials
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			fmt.Println(err)
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
		token, refreshToken, err := Login(credentials, store)
		if err != nil {
			responder.CastAndSend(err, w)
			return
		}

		// Set http-only cookie refreshToken
		c := refreshCookie(refreshToken)
		http.SetCookie(w, &c)

		responder.Res{
			Message: "logged in",
			Data: tokenResponse{
				Token: token,
			},
		}.Send(w)
	})
}

func refresh(store Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refreshToken")
		if err != nil {
			responder.Err{
				Code: http.StatusUnauthorized,
				Err:  errors.New("no refreshToken present"),
			}.Send(w)
			return
		}

		token, refreshToken, err := Refresh(cookie.Value, store)
		if err != nil {
			responder.CastAndSend(err, w)
			return
		}

		// set http-only cookie refreshToken
		c := refreshCookie(refreshToken)
		http.SetCookie(w, &c)

		responder.Res{
			Message: "token refreshed",
			Data: tokenResponse{
				Token: token,
			},
		}.Send(w)
	})
}

// me sends the user making the request
func me(w http.ResponseWriter, r *http.Request) {
	user, err := GetUser(r)
	if err != nil {
		responder.CastAndSend(err, w)
		return
	}

	responder.Res{
		Data: user,
	}.Send(w)
}

func logout(w http.ResponseWriter, r *http.Request) {
	c := refreshCookie("")
	http.SetCookie(w, &c)
	responder.Res{
		Message: "logged out",
	}.Send(w)
}
