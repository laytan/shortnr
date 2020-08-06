package click

import (
	"database/sql"
)

// Storage is the storage interface for clicks
type Storage interface {
	Create(click Click) (uint, error)
	GetAll(linkID string) []Click
}

// MysqlStorage implements Storage as click storage
type MysqlStorage struct {
	Conn *sql.DB
}

// Create adds a new click to the db and returns its id
func (m MysqlStorage) Create(click Click) (uint, error) {
	res, err := m.Conn.Exec("INSERT INTO clicks (link_id, createdAt) VALUES (?, ?)", click.LinkID, click.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

// GetAll returns all clicks of a link
func (m MysqlStorage) GetAll(linkID string) []Click {
	rows, err := m.Conn.Query("SELECT * FROM clicks WHERE link_id = ?", linkID)
	if err != nil {
		return make([]Click, 0)
	}

	clicks := make([]Click, 0)
	for rows.Next() {
		var click Click
		rows.Scan(&click.ID, &click.LinkID, &click.CreatedAt)
		clicks = append(clicks, click)
	}

	return clicks
}

// MemoryStorage implements Storage as in memory storage, used for fast testing
type MemoryStorage struct {
	Clicks []Click
}

// Create adds a click to the internal slice
func (m *MemoryStorage) Create(click Click) (uint, error) {
	click.ID = uint(len(m.Clicks) - 1)
	m.Clicks = append(m.Clicks, click)
	return click.ID, nil
}

// GetAll returns all clicks for a specific link
func (m MemoryStorage) GetAll(linkID string) []Click {
	clicks := make([]Click, 0)
	for _, click := range m.Clicks {
		if click.LinkID == linkID {
			clicks = append(clicks, click)
		}
	}
	return clicks
}
