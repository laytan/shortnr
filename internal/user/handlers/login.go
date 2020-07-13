package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/laytan/shortnr/internal/user/storage"
	"golang.org/x/crypto/bcrypt"
)

type loginBody struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}

// Login checks credentials and issues a jwt token
func Login(store storage.Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse body
		var body loginBody
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

		// Check if the user exists
		user, exists := store.GetByEmail(body.Email)
		if !exists {
			w.WriteHeader(422)
			w.Write([]byte(`{"message": "Invalid credentials"}`))
			return
		}

		// Check the password matches
		if notMatchErr := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(body.Password)); notMatchErr != nil {
			w.WriteHeader(422)
			w.Write([]byte(`{"message": "Invalid credentials"}`))
		}

		// Sign JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":        user.ID,
			"email":     user.Email,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		})

		secret, exists := os.LookupEnv("JWT_SECRET")
		if !exists {
			w.WriteHeader(500)
			panic("JWT_SECRET not set")
		}

		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			fmt.Printf("Error signing JWT: %v \n", err)
			w.WriteHeader(500)
			return
		}

		// Send token back
		w.Write([]byte(fmt.Sprintf(`{"message": "Logged in!", "token": %q}`, tokenString)))
	})
}
