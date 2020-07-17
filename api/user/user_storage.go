package user

import (
	"database/sql"
	"fmt"
)

// Storage interacts with the user store
type Storage interface {
	Get(id uint) (User, bool)
	GetByEmail(email string) (User, bool)
	Set(user User) bool
}

// Memory stores users in memory
type MemoryStorage struct {
	Users []User
}

// Set adds the user to the store
func (m *MemoryStorage) Set(user User) bool {
	m.Users = append(m.Users, user)
	return true
}

// Get returns the user with the given id and if it exists
func (m MemoryStorage) Get(id uint) (User, bool) {
	for _, user := range m.Users {
		if user.ID == id {
			return user, true
		}
	}
	return User{}, false
}

// GetByEmail returns the user with the given email and if it exists
func (m MemoryStorage) GetByEmail(email string) (User, bool) {
	for _, user := range m.Users {
		if user.Email == email {
			return user, true
		}
	}
	return User{}, false
}

// Mysql stores users in mysql database
type MysqlStorage struct {
	Conn *sql.DB
}

// Set a new user in the database
func (m MysqlStorage) Set(user User) bool {
	_, err := m.Conn.Exec(`INSERT INTO users (email, hash, createdAt) VALUES(?, ?, ?)`, user.Email, user.Hash, user.CreatedAt)
	if err != nil {
		fmt.Printf("Error in set user: %+v", err)
		return false
	}
	return true
}

// Get a user from the database with the id
func (m MysqlStorage) Get(id uint) (User, bool) {
	var user User
	err := m.Conn.QueryRow(`SELECT id, email, hash, createdAt, updatedAt FROM users WHERE id = ?`, id).Scan(&user.ID, &user.Email, &user.Hash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		fmt.Printf("Error in get users: %v", err)
		return user, false
	}
	return user, true
}

// GetByEmail returns the user associated with the email and if it exists
func (m MysqlStorage) GetByEmail(email string) (User, bool) {
	var user User
	err := m.Conn.QueryRow(`SELECT id, email, hash, createdAt, updatedAt FROM users where email = ?`, email).Scan(&user.ID, &user.Email, &user.Hash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, false
	}
	return user, true
}
