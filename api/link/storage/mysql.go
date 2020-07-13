package storage

import (
	"database/sql"
	"fmt"

	// Need to import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type link struct {
	ID  string
	URL string
}

// Mysql is a storage implementation using a mysql database connection
type Mysql struct {
	Conn *sql.DB
}

// Get returns the url associated with the id given and if it exists
func (d Mysql) Get(id string) (string, bool) {
	var link link
	err := d.Conn.QueryRow("SELECT id, url FROM links WHERE id = ?", id).Scan(&link.ID, &link.URL)
	if err != nil {
		fmt.Printf("Error in Get: %+v", err)
		return "", false
	}

	return link.URL, true
}

// Set inserts a new row to the database defining a link
func (d Mysql) Set(id string, url string) {
	_, err := d.Conn.Query("INSERT INTO links VALUES (?, ?)", id, url)
	if err != nil {
		fmt.Printf("Error in Set: %+v", err)
	}
}

// Contains checks if the database has the specified url in it
func (d Mysql) Contains(url string) bool {
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
