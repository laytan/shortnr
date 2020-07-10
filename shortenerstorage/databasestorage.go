package shortenerstorage

import (
	"database/sql"
	"fmt"
	"os"

	// Need to import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// NewDatabaseStorage returns a DatabaseStorage struct with a database connection from the environment variables
func NewDatabaseStorage() DatabaseStorage {
	// Get connection vars out of env
	username, exists := os.LookupEnv("DB_USERNAME")
	if !exists {
		panic("No DB_USERNAME env variable set")
	}
	password, exists := os.LookupEnv("DB_PASSWORD")
	if !exists {
		panic("No DB_PASSWORD env variable set")
	}
	database, exists := os.LookupEnv("DB_DATABASE")
	if !exists {
		panic("No DB_DATABASE env variable set")
	}

	// Connect to mysql
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", username, password, database))
	if err != nil {
		panic(err.Error())
	}

	// Return struct with connection
	return DatabaseStorage{Conn: db}
}

type link struct {
	ID  string
	URL string
}

// DatabaseStorage is a storage implementation using a mysql database connection
type DatabaseStorage struct {
	Conn *sql.DB
}

// Get returns the url associated with the id given and if it exists
func (d DatabaseStorage) Get(id string) (string, bool) {
	var link link
	err := d.Conn.QueryRow("SELECT id, url FROM links WHERE id = ?", id).Scan(&link.ID, &link.URL)
	if err != nil {
		fmt.Printf("Error in Get: %+v", err)
		return "", false
	}

	return link.URL, true
}

// Set inserts a new row to the database defining a link
func (d DatabaseStorage) Set(id string, url string) {
	_, err := d.Conn.Query("INSERT INTO links VALUES (?, ?)", id, url)
	if err != nil {
		fmt.Printf("Error in Set: %+v", err)
	}
}

// Contains checks if the database has the specified url in it
func (d DatabaseStorage) Contains(url string) bool {
	results, err := d.Conn.Query("SELECT COUNT(id) FROM links WHERE url = ?", url)
	if err != nil {
		fmt.Printf("Error in Contains: %+v", err)
		return false
	}

	count := 0
	for results.Next() {
		err := results.Scan(&count)
		if err != nil {
			fmt.Printf("Error in Contains: %+v", err)
			count = 0
		}
	}

	return count > 0
}
