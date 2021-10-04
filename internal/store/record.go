package store

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

var validate = validator.New()

var ErrUnparsable = errors.New("could not parse input into a record")
var ErrInvalidFields = errors.New("one or more fields is missing or invalid")

type record struct {
	Title   string `yaml:"title" validate:"required"`
	Version string `yaml:"version" validate:"required"`
	// dive tag option is necessary to validate fields in the nested struct.
	Maintainers []maintainer `yaml:"maintainers" validate:"required,gt=0,dive"`
	Company     string       `yaml:"company" validate:"required"`
	Website     string       `yaml:"website" validate:"required,url"`
	Source      string       `yaml:"source" validate:"required,url"`
	License     string       `yaml:"license" validate:"required"`
	Description string       `yaml:"description" validate:"required"`
}

type maintainer struct {
	Name  string `yaml:"name" validate:"required"`
	Email string `yaml:"email" validate:"required,email"`
}

// newRecord creates a new record from a raw stream of bytes.
// Returns an error if either the stream is unparsable or the created rawRecord doesn't conform to the schema.
func newRecord(rawRecord []byte) (*record, error) {
	var r record
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
