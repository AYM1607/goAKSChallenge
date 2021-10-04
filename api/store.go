package api

type MetaStore interface {
	Append([]byte) error
	Search([]SearchTerm) []MetaRecord
}
