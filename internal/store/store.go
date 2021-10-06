package store

import (
	"errors"
	"sync"

	"github.com/AYM1607/goAKSChallenge/api"
)

// TODO: create an index to have fast search for fields.

var ErrUnparsable = errors.New("could not parse input into a record")

type Store struct {
	// Use a read/write mutex to allow performant concurrent reads.
	mu      sync.RWMutex
	indexes map[api.SearchField]storeIndex
}

func New() (*Store, error) {
	// Create indexes for every possible search field.
	indexes := map[api.SearchField]storeIndex{}
	for _, searchField := range api.ValidSearchFieldValues() {
		isFullText := false
		if searchField == api.SearchFieldDescription {
			isFullText = true
		}
		index, err := newIndex(isFullText)
		// If any of the indexes failes to be initialized the store won't work
		// correctly and thus we should abort the whole operation.
		if err != nil {
			return nil, err
		}
		indexes[searchField] = index
	}

	return &Store{indexes: indexes}, nil
}

func (s *Store) Append(rawRecord []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, err := newRecord(rawRecord)
	if err != nil {
		return err
	}

	for _, field := range api.ValidSearchFieldValues() {
		if field == api.SearchFieldMaintainerEmail ||
			field == api.SearchFieldMaintainerName {
			continue
		}
		fieldValue, err := record.FieldValueFromSearchField(field)
		// This should not happend because we're skipping the invalid fields
		// but I'm leaving it as a safeguard.
		if err != nil {
			return err
		}
		// If we fail to add the record to any index searches won't work correclty so abort the whole operation.
		err = s.indexes[field].Index(record, fieldValue)
		if err != nil {
			return err
		}
	}

	// Since maintainers is a list it needs a separate implementation.
	for _, maintainer := range record.Maintainers {
		// If we fail to add the record to the maintainer fields indexes searches won't work correctly so we should abort.
		err := s.indexes[api.SearchFieldMaintainerEmail].Index(record, maintainer.Email)
		if err != nil {
			return err
		}
		err = s.indexes[api.SearchFieldMaintainerName].Index(record, maintainer.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) Search(joinMethod api.SearchJoinMethod,
	terms []api.SearchTerm) ([]*api.MetaRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(terms) == 0 {
		return nil, errors.New("at least one search term must be provided")
	}

	results := []*api.MetaRecord{}
	// TODO: deduplicate for "or" join method and ensure that a result is returned
	// for each of the terms in the case of the "and" join method.
	for _, term := range terms {
		// TODO: this searches should be performed concurrently.
		// Skipping indexes for fields other than description for now.
		index := s.indexes[term.Field]
		indexResults, err := index.Search(term.Query)
		// NOTE: We could argue on whether an error from a single index should fail the entire operation.
		// For the purposes of this challenge I'll allow the operation to proceed just in aces other
		// indices are able to return valid results.
		if err == nil {
			results = append(results, indexResults...)
		}
	}
	return results, nil
}
