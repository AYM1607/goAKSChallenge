package store

import (
	"errors"
	"fmt"
	"goAKSChallenge/api"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

var validate = validator.New()

var ErrUnparsable = errors.New("could not parse input into a record")

type Store struct {
	mu      sync.RWMutex
	records []api.MetaRecord
}

// newRecord creates a new record from a raw stream of bytes.
// Returns an error if either the stream is unparsable or the created rawRecord doesn't conform to the schema.
func newRecord(rawRecord []byte) (*api.MetaRecord, error) {
	var r api.MetaRecord
	if err := yaml.Unmarshal(rawRecord, &r); err != nil {
		return nil, ErrUnparsable
	}

	err := validate.Struct(r)
	if err != nil {
		fieldsWithErrors := schemaErrorFields(err.(validator.ValidationErrors))
		errString := fmt.Sprintf("the following field(s) are missing or invalid: %s",
			strings.Join(fieldsWithErrors, ","))
		return nil, errors.New(errString)
	}

	return &r, nil
}

// schemaErrorFields extracts and returns the field names that caused errors
// during schema validation from a validator.ValidationErrors.
func schemaErrorFields(errors validator.ValidationErrors) []string {
	res := []string{}
	for _, err := range errors {
		res = append(res, err.Field())
	}
	return res
}
