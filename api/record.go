package api

type Record struct {
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
