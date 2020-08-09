package link

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/laytan/shortnr/api/user"
	"github.com/laytan/shortnr/pkg/responder"
	"github.com/rs/xid"
)

// Create adds a new shortened link to the store
func Create(link Link, store Storage) (Link, error) {
	validationErr := validator.New().Struct(link)
	if validationErr != nil {
		return Link{}, responder.Err{
			Code: http.StatusBadRequest,
			Err:  errors.New("validation failed"),
		}
	}

	// Generate a unique id if no id is present
	if len(link.ID) == 0 {
		// Generate a unique ID
		id := xid.New().String()
		link.ID = id
	}

	link.CreatedAt = time.Now().Format(time.RFC3339)

	// Store in our storage
	store.Create(link)
	return link, nil
}

// Destroy removes a link if the requester made it and the link exists
func Destroy(linkID string, requester user.User, store Storage) error {
	link, exists := store.Get(linkID)
	if !exists {
		return responder.Err{
			Code: http.StatusNotFound,
			Err:  errors.New("link does not exist"),
		}
	}

	// Check if user is link creator
	if link.UserID != requester.ID {
		return responder.Err{
			Code: http.StatusUnauthorized,
			Err:  errors.New("this link is not created by you"),
		}
	}

	// Delete link
	deleted := store.Delete(linkID)
	if !deleted {
		return responder.Err{
			Code: http.StatusInternalServerError,
			Err:  errors.New("link could not be deleted"),
		}
	}

	return nil
}
