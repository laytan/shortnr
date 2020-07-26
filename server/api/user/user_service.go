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
func Login(credentials Credentials, store Storage) (string, string, error) {
	// validate credentials
	validationErr := validator.New().Struct(credentials)
	if validationErr != nil {
		return "", "", responder.Err{
			Code: http.StatusBadRequest,
			Err:  errors.New("credentials validation failed"),
		}
	}

	// Check if the user exists
	user, exists := store.GetByEmail(credentials.Email)
	if !exists {
		return "", "", responder.Err{
			Code: http.StatusUnauthorized,
			Err:  errors.New("invalid credentials"),
		}
	}

	// Check the password matches
	if notMatchErr := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(credentials.Password)); notMatchErr != nil {
		return "", "", responder.Err{
			Code: http.StatusUnauthorized,
			Err:  errors.New("invalid credentials"),
		}
	}

	token, refreshToken, err := SignUserToken(user)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

// SignUserToken creates a jwt token for the given user, errors are of type responder.Err
func SignUserToken(user User) (string, string, error) {
	// Sign JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        user.ID,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
		// Expire in 15 minutes
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})

	// Generate refresh token
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secret, exists := os.LookupEnv("JWT_SECRET")
	if !exists {
		return "", "", responder.Err{
			Code: http.StatusInternalServerError,
			Err:  errors.New("JWT_SECRET environment variable not set"),
		}
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", responder.Err{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("error signing JWT: %+v", err),
		}
	}

	refreshString, err := refresh.SignedString([]byte(secret))
	if err != nil {
		return "", "", responder.Err{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("error signing JWT: %+v", err),
		}
	}

	return tokenString, refreshString, nil
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

// Refresh generates new tokens based on the given refresh token
func Refresh(refreshToken string, store Storage) (string, string, error) {
	if len(refreshToken) < 2 {
		return "", "", responder.Err{
			Code: http.StatusBadRequest,
			Err:  errors.New("no refresh token provided"),
		}
	}

	// Parse the token
	parsedToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		secret, exists := os.LookupEnv("JWT_SECRET")
		if !exists {
			return nil, responder.Err{
				Code: http.StatusInternalServerError,
				Err:  errors.New("no jwt secret set in env"),
			}
		}

		return []byte(secret), nil
	})
	if err != nil {
		return "", "", err
	}

	// Parse claims for user id
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return "", "", responder.Err{
			Code: http.StatusUnauthorized,
			Err:  errors.New("invalid token"),
		}
	}

	// check if the user exists and is allowed to login
	userID := uint(claims["id"].(float64))
	user, exists := store.Get(userID)
	if !exists {
		return "", "", responder.Err{
			Code: http.StatusUnauthorized,
			Err:  errors.New("user does not exist anymore"),
		}
	}

	return SignUserToken(user)
}
