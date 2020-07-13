package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/laytan/shortnr/internal/user/storage"
)

type contextKey int

// AuthenticatedUserKey is the key on the request context for the authenticated user
const AuthenticatedUserKey contextKey = 0

// JwtAuthorization checks for a jwt token and sets the user on the request
func JwtAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Make sure there is a token to substring next
		if len(authHeader) < 10 {
			w.WriteHeader(403)
			return
		}

		// Take out actual token, removes Bearer
		token := string([]rune(authHeader)[7:])

		// Parse the token
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			w.WriteHeader(403)
			return
		}

		// Make sure the token is valid
		if _, ok := parsedToken.Claims.(jwt.Claims); !ok || !parsedToken.Valid {
			w.WriteHeader(403)
			return
		}

		// Extract user from jwt token
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok || !parsedToken.Valid {
			w.WriteHeader(403)
			return
		}

		// Turn map of claims into json
		jsonClaims, err := json.Marshal(claims)
		if err != nil {
			w.WriteHeader(403)
			return
		}

		// Turn json into User type
		user := storage.User{}
		if err := json.Unmarshal(jsonClaims, &user); err != nil {
			w.WriteHeader(403)
			return
		}

		// Store the user into the context of the request
		ctxWithUser := context.WithValue(r.Context(), AuthenticatedUserKey, user)
		rWithUser := r.WithContext(ctxWithUser)

		next.ServeHTTP(w, rWithUser)
	})
}
