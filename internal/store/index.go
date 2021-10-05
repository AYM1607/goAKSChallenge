package store

import (
	"github.com/AYM1607/goAKSChallenge/api"
	"github.com/blevesearch/bleve/v2"
)

type storeIndex interface {
	Index(*api.MetaRecord, string) error
	Search(string) ([]*api.MetaRecord, error)
}

func newIndex(name string, isFullText bool) (storeIndex, error) {
	if isFullText {
		// Create bleve index.
		mapping := bleve.NewIndexMapping()
		bleveIndex, err := bleve.NewMemOnly(mapping)
		if err != nil {
			return nil, err
		}

		index := fullTextSearchIndex{name: name, bleveIndex: bleveIndex}
		return &index, nil
	}
	// TODO: return the implementation of the ExactMatchSearchIndex
	return nil, nil

}

// NOTE: Implementig a full text search would have been too much work for the purposes
// of this challenge but I still wanted to have the feature available for the description field.
// The bleve library is probably too overkill but this purpose, but once again, I just wanted
// to add the feature regardless of the size of the final binary. Further optimizations
// could be possible if we narrowed the requirements for the search.
type fullTextSearchIndex struct {
	name       string
	bleveIndex bleve.Index
	mapping    map[string]*api.MetaRecord
}

func (i fullTextSearchIndex) Index(*api.MetaRecord, string) error {
	// TODO: since bleve only accepts strings as identifiers we have to have an
	// internal mapping from a UID (can be ULID) to the pointer intself.
	return nil
}

func (i fullTextSearchIndex) Search(string) ([]*api.MetaRecord, error) {
	return nil, nil
}

// Implement an exact match index with a trie.
type exactMatchSearchIndex struct {
}
