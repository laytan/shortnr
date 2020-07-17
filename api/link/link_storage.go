package link

import (
	"database/sql"
	"fmt"
)

// Storage defines the methods required to interact with a link storing method
type Storage interface {
	Get(id string) (string, bool)
	Set(id string, url string)
	Contains(url string) bool
}

type link struct {
	ID  string
	URL string
}

// Mysql is a storage implementation using a mysql database connection
type MysqlStorage struct {
	Conn *sql.DB
}

// Get returns the url associated with the id given and if it exists
func (d MysqlStorage) Get(id string) (string, bool) {
	var link link
	err := d.Conn.QueryRow("SELECT id, url FROM links WHERE id = ?", id).Scan(&link.ID, &link.URL)
	if err != nil {
		fmt.Printf("Error in Get: %+v", err)
		return "", false
	}

	return link.URL, true
}

// Set inserts a new row to the database defining a link
func (d MysqlStorage) Set(id string, url string) {
	_, err := d.Conn.Query("INSERT INTO links VALUES (?, ?)", id, url)
	if err != nil {
		fmt.Printf("Error in Set: %+v", err)
	}
}

// Contains checks if the database has the specified url in it
func (d MysqlStorage) Contains(url string) bool {
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

type MemoryStorage struct {
	InternalMap map[string]string
}

// Get returns the url associated with the given id
func (m MemoryStorage) Get(id string) (string, bool) {
	redirectURL, ok := m.InternalMap[id]
	return redirectURL, ok
}

// Set adds the given id to url to the map
func (m MemoryStorage) Set(id string, url string) {
	m.InternalMap[id] = url
}

// Contains checks if the map has the url in it
func (m MemoryStorage) Contains(url string) bool {
	contains := false
	for _, v := range m.InternalMap {
		if v == url {
			contains = true
			break
		}
	}
	return contains
}
