package routes

import (
	"net/http"

	"github.com/laytan/shortnr/internal/user/middlewares"

	"github.com/gorilla/mux"
	"github.com/laytan/shortnr/internal/user/handlers"
	"github.com/laytan/shortnr/internal/user/storage"
)

// Set adds the routes needed for users service
func Set(r *mux.Router, storage storage.Storage) {
	r.HandleFunc("/signup", handlers.Signup(storage)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.Login(storage)).Methods(http.MethodPost)

	authR := r.PathPrefix("").Subrouter()
	authR.Use(middlewares.JwtAuthorization)
	authR.HandleFunc("/me", handlers.Me()).Methods(http.MethodGet)
}
