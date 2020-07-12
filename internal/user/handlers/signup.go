package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/laytan/shortnr/internal/user/storage"
	"golang.org/x/crypto/bcrypt"
)

type signupBody struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}

// Signup tries to create a new user into the database
func Signup(store storage.Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse body
		var body signupBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(`{"message": "Unable to parse JSON body"}`))
			return
		}

		// validate email and password
		validationErr := validator.New().Struct(body)
		if validationErr != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf(`{"message": "Validation failed", "error": %q}`, validationErr)))
			return
		}

		// check not already in use
		_, exists := store.GetByEmail(body.Email)
		if exists {
			w.WriteHeader(409)
			w.Write([]byte(`{"message": "Email in use"}`))
			return
		}

		// hash password
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
		if hashErr != nil {
			w.WriteHeader(500)
			fmt.Printf("Error while hashing password: %v", hashErr)
			w.Write([]byte(`{"message": "Error hashing password"}`))
			return
		}

		// insert user into storage
		user := storage.User{Email: body.Email, Hash: string(hash), CreatedAt: time.Now()}
		inserted := store.Set(user)
		if !inserted {
			w.WriteHeader(500)
			w.Write([]byte(`{"message": "Database error"}`))
			return
		}

		// respond with success message
		w.Write([]byte(`{"message": "Signed up!"}`))
	})
}
