package api

import "errors"

var ErrFieldLookupNotSupported = errors.New("the lookup of this search field is not supported")

type MetaRecord struct {
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

// fieldValueFromSearchField returns the struct field value from a given SearchField.
// NOTE: The reflection api could probably be used to make this implementations simpler,
// but I'm pretty sure that's not the purpose of that api.
// This implementation is flaky because any new value of SearchField means this
// function has to be updated and there's no mechanism to produce a compilation
// error if we forget to add it.
func (r *MetaRecord) FieldValueFromSearchField(field SearchField) (string, error) {

	switch field {
	case SearchFieldCompany:
		return r.Company, nil
	case SearchFieldLicense:
		return r.License, nil
	case SearchFieldMaintainerEmail:
		return "", ErrFieldLookupNotSupported
	case SearchFieldMaintainerName:
		return "", ErrFieldLookupNotSupported
	case SearchFieldSource:
		return r.Source, nil
	case SearchFieldTitle:
		return r.Title, nil
	case SearchFieldVersion:
		return r.Version, nil
	case SearchFieldWebsite:
		return r.Website, nil
	case SearchFieldDescription:
		return r.Description, nil
	}
	return "", nil
}
