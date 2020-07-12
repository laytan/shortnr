package storage

import (
	"database/sql"
	"fmt"
)

// Mysql stores users in mysql database
type Mysql struct {
	Conn *sql.DB
}

// Set a new user in the database
func (m Mysql) Set(user User) bool {
	fmt.Println(user)
	_, err := m.Conn.Exec(`INSERT INTO users (email, hash, createdAt) VALUES(?, ?, ?)`, user.Email, user.Hash, user.CreatedAt)
	if err != nil {
		fmt.Printf("Error in set user: %+v", err)
		return false
	}
	return true
}

// Get a user from the database with the id
func (m Mysql) Get(id uint) (User, bool) {
	var user User
	err := m.Conn.QueryRow(`SELECT id, email, hash, createdAt, updatedAt FROM users WHERE id = ?`, id).Scan(&user.ID, &user.Email, &user.Hash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		fmt.Printf("Error in get users: %v", err)
		return user, false
	}
	return user, true
}

// GetByEmail returns the user associated with the email and if it exists
func (m Mysql) GetByEmail(email string) (User, bool) {
	var user User
	err := m.Conn.QueryRow(`SELECT id, email, hash, createdAt, updatedAt FROM users where email = ?`, email).Scan(&user.ID, &user.Email, &user.Hash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, false
	}
	return user, true
}
