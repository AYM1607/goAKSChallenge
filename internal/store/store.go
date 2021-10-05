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

func (s *Store) Search(joinMethod api.SearchJoinMethod,
	terms []api.SearchTerm) ([]*api.MetaRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(terms) == 0 {
		return nil, errors.New("at least one search term must be provided")
	}

	return nil, nil
}
