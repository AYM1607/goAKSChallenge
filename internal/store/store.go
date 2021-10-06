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

	// Used to keep track of how many times a record has shown up in term searches.
	foundRecordsCounts := map[*api.MetaRecord]int{}
	results := []*api.MetaRecord{}

	// Results returned from individual index searches.
	indexesQueryResults := make(chan []*api.MetaRecord)

	// Goroutines sync.
	done := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(len(terms))

	for _, term := range terms {
		go func(term api.SearchTerm) {
			index := s.indexes[term.Field]
			indexQueryResult, err := index.Search(term.Query)
			// NOTE: We could argue on whether an error from a single index should fail the entire operation.
			// For the purposes of this challenge I'll allow the operation to proceed just in aces other
			// indexes are able to return valid results.
			if err == nil {
				indexesQueryResults <- indexQueryResult
			}
			wg.Done()
		}(term)
	}

	go func() {
		wg.Wait()
		done <- true
	}()

	// Process all the results coming from the indexes queries and cleanup when they're all done.
P:
	for {
		select {
		// Clean channels and finish the loop if there's no more index requests to process.
		case <-done:
			close(done)
			close(indexesQueryResults)
			break P
		case indexQueryResult := <-indexesQueryResults:
			for _, match := range indexQueryResult {
				switch joinMethod {
				// Only add the record to the results on its first match.
				case api.SearchJoinMethodOR:
					if foundRecordsCounts[match] == 0 {
						foundRecordsCounts[match] = 1
						results = append(results, match)
					}
				// Only add the record to the results if its been matched in all term searches.
				case api.SearchJoinMethodAND:
					foundRecordsCounts[match] += 1
					if foundRecordsCounts[match] == len(terms) {
						results = append(results, match)
					}
				}
			}
		}
	}

	return results, nil
}
