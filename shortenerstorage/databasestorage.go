package shortenerstorage

import (
	"database/sql"
	"fmt"
	"os"
)

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
	host, exists := os.LookupEnv("DB_HOST")
	if !exists {
		panic("No DB_USERNAME env variable set")
	}
	port, exists := os.LookupEnv("DB_PORT")
	if !exists {
		panic("No DB_PORT env variable set")
	}
	database, exists := os.LookupEnv("DB_DATABASE")
	if !exists {
		panic("No DB_DATABASE env variable set")
	}

	// Connect to mysql
	db, err := sql.Open("mysql", fmt.Sprintf("%q:%q@%q:%q/%q", username, password, host, port, database))
	if err != nil {
		panic(err.Error())
	}

	// Return struct with connection
	return DatabaseStorage{Conn: db}
}

type Link struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type DatabaseStorage struct {
	Conn *sql.DB
}

func (d DatabaseStorage) Get(id string) (string, bool) {
	var link Link
	err := d.Conn.QueryRow("SELECT id, url FROM links WHERE id = ?", id).Scan(&link)
	if err != nil {
		fmt.Printf("Error in Get: %+v", err)
		return "", false
	}

	return link.URL, true
}

func (d DatabaseStorage) Set(id string, url string) {
	_, err := d.Conn.Query("INSERT INTO links VALUES (?, ?)", id, url)
	if err != nil {
		fmt.Printf("Error in Set: %+v", err)
	}
}

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
