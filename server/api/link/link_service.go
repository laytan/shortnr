package link

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator"
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

	// Generate a unique ID
	id := xid.New().String()

	link.ID = id
	link.CreatedAt = time.Now().Format(time.RFC3339)

	// Store in our storage
	store.Create(link)
	return link, nil
}
