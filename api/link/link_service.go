package link

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/laytan/shortnr/pkg/responder"
	"github.com/rs/xid"
)

// Create adds a new shortened link to the store
func Create(link Link, store Storage) (string, error) {
	validationErr := validator.New().Struct(link)
	if validationErr != nil {
		return "", responder.Err{
			Code: http.StatusBadRequest,
			Err:  errors.New("validation failed"),
		}
	}

	// Generate a unique ID
	id := xid.New().String()

	// Store in our storage
	store.Set(id, link.URL)
	return id, nil
}
