package store

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/AYM1607/goAKSChallenge/api"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

// TODO: create an index to have fast search for fields.

var validate = validator.New()

var ErrUnparsable = errors.New("could not parse input into a record")

type Store struct {
	// Use a read/write mutex to allow performant concurrent reads.
	mu      sync.RWMutex
	records []*api.MetaRecord
}

func New() *Store {
	return &Store{}
}

func (s *Store) Append(rawRecord []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, err := newRecord(rawRecord)
	if err != nil {
		return err
	}

	// TODO: When we make use of indexes they'll have to be updated here.

	s.records = append(s.records, record)
	return nil
}

// TODO: Perhaps the record logic should be moved to its own file to separate concerns.

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
		// The Namespace uses the go struct names instead of the yaml tags.
		// Capitalization can be confusing but its a tradeoff to allow knowing
		// exactly what field caused the error.

		// Delete the root element of the namespace. Having the name of the internal go struct can throw off users.
		ns := err.Namespace()

		ns = ns[strings.Index(ns, ".")+1:]

		res = append(res, ns)
	}
	return res
}
