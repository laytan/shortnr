package storage

import (
	"database/sql"
	"encoding/json"
	"time"
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

// Storage interacts with the user store
type Storage interface {
	Get(id uint) (User, bool)
	GetByEmail(email string) (User, bool)
	Set(user User) bool
}
