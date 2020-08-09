package validation

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator"
)

// Custom error messages per validation tag
var parsers = map[string]func(validator.FieldError, string) string{
	"url": func(e validator.FieldError, field string) string {
		return fmt.Sprintf("%s is not a valid URL", field)
	},
	"required": func(e validator.FieldError, field string) string {
		return fmt.Sprintf("%s is required", field)
	},
	"max": func(e validator.FieldError, field string) string {
		return fmt.Sprintf("%s must be a maximum of %s characters", field, e.Param())
	},
	"min": func(e validator.FieldError, field string) string {
		return fmt.Sprintf("%s must be a minimum of %s characters", field, e.Param())
	},
	"email": func(e validator.FieldError, field string) string {
		return fmt.Sprintf("%s is not a valid email", field)
	},
}

// ParseError parses a validation error to display to users, note: only the first aliasses map is used
func ParseError(e error, aliasses ...map[string]string) error {
	errorStr := ""
	for i, err := range e.(validator.ValidationErrors) {
		// Get specific parser
		parser := parsers[err.Tag()]
		if parser == nil {
			panic("No parser for tag: " + err.Tag())
		}

		// If an alias is given use that instead of the field name
		field := err.Field()
		if len(aliasses) > 0 {
			alias := aliasses[0][field]
			if len(alias) > 0 {
				field = alias
			}
		}

		// Run through specific parser
		errorStr += parser(err, field)

		// Add , when it is not the last error
		if i+1 != len(e.(validator.ValidationErrors)) {
			errorStr += ", "
		}
	}
	return errors.New(errorStr)
}
