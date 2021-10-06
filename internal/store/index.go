package store

import (
	"math/rand"
	"time"

	"github.com/AYM1607/goAKSChallenge/api"
	"github.com/blevesearch/bleve/v2"
	"github.com/oklog/ulid/v2"
)

type storeIndex interface {
	Index(*api.MetaRecord, string) error
	Search(string) ([]*api.MetaRecord, error)
}

func newIndex(isFullText bool) (storeIndex, error) {
	if isFullText {
		// Create bleve index.
		mapping := bleve.NewIndexMapping()
		bleveIndex, err := bleve.NewMemOnly(mapping)
		if err != nil {
			return nil, err
		}

		index := fullTextSearchIndex{
			bleveIndex: bleveIndex,
			idMap:      map[string]*api.MetaRecord{},
		}

		return &index, nil
	}
	// TODO: return the implementation of the ExactMatchSearchIndex
	return nil, nil

}

// NOTE: Implementig a full text search would have been too much work for the purposes
// of this challenge but I still wanted to have the feature available for the description field.
// The bleve library is probably too overkill for this purpose, but once again, I just wanted
// to add the feature regardless of the size of the final binary. Further optimizations
// could be possible if we narrowed the requirements for the search.
type fullTextSearchIndex struct {
	name       string
	bleveIndex bleve.Index
	idMap      map[string]*api.MetaRecord
}

func (i fullTextSearchIndex) Index(record *api.MetaRecord, data string) error {
	// Create a string parsable UID for the record.
	// This is necessary because bleve only accepts strings as document identifiers.
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	recordId, err := ulid.New(ulid.Timestamp(t), entropy)

	if err != nil {
		return err
	}

	i.idMap[recordId.String()] = record
	err = i.bleveIndex.Index(recordId.String(), data)
	if err != nil {
		return err
	}

	return nil
}

func (i fullTextSearchIndex) Search(term string) ([]*api.MetaRecord, error) {
	// Retireve the internal ids for the records from the bleve index.
	query := bleve.NewMatchQuery(term)
	search := bleve.NewSearchRequest(query)
	searchResults, err := i.bleveIndex.Search(search)
	if err != nil {
		return nil, err
	}

	resultRecords := []*api.MetaRecord{}

	// Convert from bleve ids to MetaRecord pointers.
	// NOTE: This is linear and could perform poorly when the number of results is high.
	for _, match := range searchResults.Hits {
		resultRecords = append(resultRecords, i.idMap[match.ID])
	}
	return resultRecords, nil
}

// Implement an exact match index with a trie.
type exactMatchSearchIndex struct {
}
