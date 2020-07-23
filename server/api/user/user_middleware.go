package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/laytan/shortnr/pkg/responder"
)

type contextKey int

// AuthenticatedUserKey is the key on the request context for the authenticated user
const AuthenticatedUserKey contextKey = 0

// JwtAuthorization checks for a jwt token and sets the user on the request
func JwtAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		unauthorized := responder.Err{
			Code: http.StatusUnauthorized,
			Err:  errors.New("jwt token invalid"),
		}

		// Make sure there is a token to substring next
		if len(authHeader) < 10 {
			unauthorized.Send(w)
			return
		}

		// Take out actual token, removes Bearer
		token := string([]rune(authHeader)[7:])

		// Parse the token
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			secret, exists := os.LookupEnv("JWT_SECRET")
			if !exists {
				noEnvErr := responder.Err{
					Code: http.StatusInternalServerError,
					Err:  errors.New("no jwt secret set in env"),
				}
				noEnvErr.Send(w)
				panic(noEnvErr)
			}

			return []byte(secret), nil
		})
		if err != nil {
			unauthorized.Send(w)
			return
		}

		// Make sure the token is valid
		if _, ok := parsedToken.Claims.(jwt.Claims); !ok || !parsedToken.Valid {
			unauthorized.Send(w)
			return
		}

		// Extract user from jwt token
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok || !parsedToken.Valid {
			unauthorized.Send(w)
			return
		}

		// Turn map of claims into json
		jsonClaims, err := json.Marshal(claims)
		if err != nil {
			unauthorized.Send(w)
			return
		}

		// Turn json into User type
		user := User{}
		if err := json.Unmarshal(jsonClaims, &user); err != nil {
			unauthorized.Send(w)
			return
		}

		// Store the user into the context of the request
		ctxWithUser := context.WithValue(r.Context(), AuthenticatedUserKey, user)
		rWithUser := r.WithContext(ctxWithUser)

		next.ServeHTTP(w, rWithUser)
	})
}
