package api

import "errors"

type SearchField string

type SearchJoinMethod string

const (
	// Field enum values.
	SearchFieldTitle           = "title"
	SearchFieldVersion         = "version"
	SearchFieldMaintainerEmail = "maintainerEmail"
	SearchFieldMaintainerName  = "maintainerName"
	SearchFieldCompany         = "company"
	SearchFieldWebsite         = "website"
	SearchFieldSource          = "source"
	SearchFieldLicense         = "license"

	// Join method enum values.
	SearchJoinMethodAND = "and"
	SearchJoinMethodOR  = "or"
)

// With no native enums in Go the following 2 functions are decent validation methods.
// TODO: Consider alternative ways of parsing and validating search parameters.

// IsValid determines if the instace of SearchField is one of the valid enum values.
func (f SearchField) IsValid() error {
	switch f {
	case SearchFieldCompany,
		SearchFieldLicense,
		SearchFieldMaintainerEmail,
		SearchFieldMaintainerName,
		SearchFieldSource,
		SearchFieldTitle,
		SearchFieldVersion,
		SearchFieldWebsite:
		return nil
	}
	return errors.New("invalid search field type")
}

// IsValid determines if the instance of SearchJoinMethod is one of the valid enum values.
func (jm SearchJoinMethod) IsValid() error {
	switch jm {
	case SearchJoinMethodAND, SearchJoinMethodOR:
		return nil
	}
	return errors.New("invalid join method type")
}

type SearchTerm struct {
	Field SearchField `json:"field"`
	Query string      `json:"query"`
}
