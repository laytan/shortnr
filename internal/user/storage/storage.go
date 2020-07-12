package storage

import (
	"database/sql"
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

// Storage interacts with the user store
type Storage interface {
	Get(id uint) (User, bool)
	GetByEmail(email string) (User, bool)
	Set(user User) bool
}
