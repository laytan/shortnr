package click

import (
	"net/http"
	"time"

	"github.com/laytan/shortnr/pkg/responder"
)

// Create creates a new click
func Create(linkID string, store Storage) (Click, error) {
	click := Click{
		LinkID:    linkID,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	id, err := store.Create(click)
	if err != nil {
		return click, responder.Err{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	click.ID = id
	return click, nil
}

// GetAll returns all the clicks of a specific link
func GetAll(linkID string, store Storage) []Click {
	return store.GetAll(linkID)
}
