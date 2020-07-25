package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/laytan/shortnr/pkg/responder"
)

// User is a user
type User struct {
	ID        uint
	Email     string
	Hash      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// MarshalJSON removes hash, formats to start with lowercase, sets updatedAt to nil if invalid
func (u User) MarshalJSON() ([]byte, error) {
	user := make(map[string]interface{})
	user["id"] = u.ID
	user["email"] = u.Email
	user["createdAt"] = u.CreatedAt
	user["updatedAt"] = u.UpdatedAt.Time
	if !u.UpdatedAt.Valid {
		user["updatedAt"] = nil
	}
	return json.Marshal(user)
}

// GetUser returns the user that made this request (if logged in) or errors
func GetUser(r *http.Request) (User, error) {
	user, ok := r.Context().Value(AuthenticatedUserKey).(User)
	if !ok {
		return User{}, responder.Err{
			Code: http.StatusUnauthorized,
			Err:  errors.New("user is not logged in"),
		}
	}
	return user, nil
}

// Credentials are used on login and signup requests
type Credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
