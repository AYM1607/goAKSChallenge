package store

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

var validate = validator.New()

var ErrUnparsable = errors.New("could not parse input into a record")
var ErrInvalidFields = errors.New("one or more fields is missing or invalid")

type record struct {
	Title       string       `yaml:"title" validate:"required"`
	Version     string       `yaml:"version" validate:"required"`
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

func newRecord(rawRecord []byte) (*record, error) {
	var r record
	if err := yaml.Unmarshal(rawRecord, &r); err != nil {
		return nil, ErrUnparsable
	}

	fmt.Printf("%+v\n", r)

	err := validate.Struct(r)
	if err != nil {
		return nil, ErrInvalidFields
	}

	return &r, nil
}
