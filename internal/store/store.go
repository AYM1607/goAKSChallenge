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
	records []*api.MetaRecord
	indexes map[api.SearchField]storeIndex
}

func New() (*Store, error) {
	indexes := map[api.SearchField]storeIndex{}
	// TODO: Create the indexes for the rest of the fields.
	index, err := newIndex(true)
	if err != nil {
		return nil, err
	}
	indexes[api.SearchFieldDescription] = index
	return &Store{indexes: indexes}, err
}

func (s *Store) Append(rawRecord []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, err := newRecord(rawRecord)
	if err != nil {
		return err
	}

	// TODO: When we make use of indexes they'll have to be updated here.
	s.indexes[api.SearchFieldDescription].Index(record, record.Description)

	s.records = append(s.records, record)
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
	// TODO: query all the indexes.
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
