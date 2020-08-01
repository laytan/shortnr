package link

import (
	"database/sql"
	"fmt"
)

// Storage defines the methods required to interact with a link storing method
type Storage interface {
	Get(id string) (Link, bool)
	GetLinksFromUser(userID uint) []Link
	Create(link Link)
	Contains(url string) bool
	Delete(id string) bool
}

// type link struct {
// 	ID  string
// 	URL string
// }

// MysqlStorage is a storage implementation using a mysql database connection
type MysqlStorage struct {
	Conn *sql.DB
}

// Get returns the url associated with the id given and if it exists
func (d MysqlStorage) Get(id string) (Link, bool) {
	var link Link
	err := d.Conn.QueryRow("SELECT * FROM links WHERE id = ?", id).Scan(&link.ID, &link.URL, &link.UserID, &link.CreatedAt)
	if err != nil {
		return Link{}, false
	}

	return link, true
}

// GetLinksFromUser returns all the links created by the given userID
func (d MysqlStorage) GetLinksFromUser(userID uint) []Link {
	rows, err := d.Conn.Query("SELECT * FROM links WHERE user_id = ?", userID)
	if err != nil {
		return make([]Link, 0)
	}

	links := make([]Link, 0)
	for rows.Next() {
		var link Link
		rows.Scan(&link.ID, &link.URL, &link.UserID, &link.CreatedAt)
		links = append(links, link)
	}

	return links
}

// Create inserts a new row to the database defining a link
func (d MysqlStorage) Create(link Link) {
	_, err := d.Conn.Exec("INSERT INTO links VALUES (?, ?, ?, ?)", link.ID, link.URL, link.UserID, link.CreatedAt)
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

// Delete removes the link with the given id
func (d MysqlStorage) Delete(id string) bool {
	res, err := d.Conn.Exec("DELETE FROM links WHERE id = ?", id)
	if err != nil {
		return false
	}

	affected, err := res.RowsAffected()
	if affected == 0 || err != nil {
		return false
	}

	return true
}

// MemoryStorage stores links in memory
type MemoryStorage struct {
	Links []Link
}

// Get returns the url associated with the given id
func (m MemoryStorage) Get(id string) (Link, bool) {
	for _, link := range m.Links {
		if link.ID == id {
			return link, true
		}
	}
	return Link{}, false
}

// GetLinksFromUser gets all the links that the user with the given id has created
func (m MemoryStorage) GetLinksFromUser(userID uint) []Link {
	links := make([]Link, 0)
	for _, link := range m.Links {
		if link.UserID == userID {
			links = append(links, link)
		}
	}
	return links
}

// Create adds the given id to url to the map
func (m *MemoryStorage) Create(link Link) {
	m.Links = append(m.Links, link)
}

// Contains checks if the map has the url in it
func (m MemoryStorage) Contains(url string) bool {
	for _, v := range m.Links {
		if v.URL == url {
			return true
		}
	}
	return false
}

// Delete removes the link with the given id
func (m *MemoryStorage) Delete(id string) bool {
	// find link index
	linkIndex := -1
	for i, v := range m.Links {
		if v.ID == id {
			linkIndex = i
		}
	}
	if linkIndex == -1 {
		return false
	}

	// Remove index from slice
	m.Links[linkIndex] = m.Links[len(m.Links)-1]
	m.Links = m.Links[:len(m.Links)-1]
	return true
}
