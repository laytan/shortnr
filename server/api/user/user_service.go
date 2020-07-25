package user

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/laytan/shortnr/pkg/responder"
	"golang.org/x/crypto/bcrypt"
)

// Login handles logging in and returning a user token, erros are of type responder.Err
func Login(credentials Credentials, store Storage) (User, string, error) {
	// validate credentials
	validationErr := validator.New().Struct(credentials)
	if validationErr != nil {
		return User{}, "", responder.Err{
			Code: http.StatusBadRequest,
			Err:  errors.New("credentials validation failed"),
		}
	}

	// Check if the user exists
	user, exists := store.GetByEmail(credentials.Email)
	if !exists {
		return User{}, "", responder.Err{
			Code: http.StatusUnauthorized,
			Err:  errors.New("invalid credentials"),
		}
	}

	// Check the password matches
	if notMatchErr := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(credentials.Password)); notMatchErr != nil {
		return User{}, "", responder.Err{
			Code: http.StatusUnauthorized,
			Err:  errors.New("invalid credentials"),
		}
	}

	token, err := SignUserToken(user)
	if err != nil {
		return User{}, "", err
	}

	return user, token, nil
}

// SignUserToken creates a jwt token for the given user, errors are of type responder.Err
func SignUserToken(user User) (string, error) {
	// Sign JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        user.ID,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})

	secret, exists := os.LookupEnv("JWT_SECRET")
	if !exists {
		return "", responder.Err{
			Code: http.StatusInternalServerError,
			Err:  errors.New("JWT_SECRET environment variable not set"),
		}
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", responder.Err{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("error signing JWT: %+v", err),
		}
	}

	return tokenString, nil
}

// Signup handles signing up a user, errors are of type responder.Err
func Signup(creds Credentials, store Storage) error {
	// validate email and password
	validationErr := validator.New().Struct(creds)
	if validationErr != nil {
		return responder.Err{
			Code: http.StatusBadRequest,
			Err:  errors.New("validation failed"),
		}
	}

	// check not already in use
	_, exists := store.GetByEmail(creds.Email)
	if exists {
		return responder.Err{
			Code: http.StatusConflict,
			Err:  errors.New("email in use"),
		}
	}

	// hash password
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(creds.Password), 14)
	if hashErr != nil {
		return responder.Err{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("error while hashing password: %+v", hashErr),
		}
	}

	// insert user into storage
	user := User{Email: creds.Email, Hash: string(hash), CreatedAt: time.Now()}
	inserted := store.Set(user)
	if !inserted {
		return responder.Err{
			Code: http.StatusInternalServerError,
			Err:  errors.New("error while inserting user into storage"),
		}
	}
	return nil
}
