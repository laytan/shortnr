package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/laytan/shortnr/internal/user/middlewares"
	"github.com/laytan/shortnr/internal/user/storage"
)

// Me returns the logged in user information
func Me() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(middlewares.AuthenticatedUserKey).(storage.User)
		jsonUser, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(403)
			return
		}

		w.WriteHeader(200)
		w.Write(jsonUser)
	})
}
