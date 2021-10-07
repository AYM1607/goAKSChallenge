package store

import (
	"testing"

	"github.com/AYM1607/goAKSChallenge/api"
	"github.com/stretchr/testify/require"
)

var records = []struct {
	record *api.MetaRecord
	data   string
}{
	{record: &api.MetaRecord{}, data: "App 1"},
	{record: &api.MetaRecord{}, data: "App 2"},
	{record: &api.MetaRecord{}, data: "hello@email.com"},
	{record: &api.MetaRecord{}, data: "This is some sample text"},
	{record: &api.MetaRecord{}, data: "description of an app"},
	{record: &api.MetaRecord{}, data: `Hello this is a very long string that will
	make sure the full text indexer is able to catch subtelties`},
}

func TestIndexing(t *testing.T) {

	data := []struct {
		name        string
		record      *api.MetaRecord
		data        string
		shouldFail  bool
		testMessage string
	}{
		{
			name:        "Both fields valid",
			record:      &api.MetaRecord{},
			data:        "some text",
			shouldFail:  false,
			testMessage: "index should index correctly if provided with valid data",
		},
		{
			name:        "Nil record pointer",
			record:      nil,
			data:        "some text",
			shouldFail:  true,
			testMessage: "index should not index if provided with a nil pointer",
		},
		{
			name:        "Zero value for data",
			record:      &api.MetaRecord{},
			data:        "",
			shouldFail:  true,
			testMessage: "index should not index if provided with the zero value for data",
		},
		{
			name:        "Both fields invalid",
			record:      nil,
			data:        "",
			shouldFail:  true,
			testMessage: "index should not index if the record pointer is nil and the data is the zero value",
		},
	}

	// Full text.
	index, err := newIndex(true)
	require.NoError(t, err)

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err = index.Index(d.record, d.data)
			if d.shouldFail {
				require.Error(t, err, d.testMessage)
			} else {
				require.NoError(t, err, d.testMessage)
			}
		})
	}

	// Exact match.
	index, err = newIndex(false)
	require.NoError(t, err)

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err = index.Index(d.record, d.data)
			if d.shouldFail {
				require.Error(t, err, d.testMessage)
			} else {
				require.NoError(t, err, d.testMessage)
			}
		})
	}
}

func TestFullTextSearching(t *testing.T) {
	index, err := newIndex(true)
	require.NoError(t, err)

	for _, r := range records {
		index.Index(r.record, r.data)
	}

	_, err = index.Search("")
	require.Error(t, err, "search should fail if provided an empty term")

	data := []struct {
		term    string
		results []*api.MetaRecord
	}{
		{term: "app", results: []*api.MetaRecord{records[0].record, records[1].record, records[4].record}},
		{term: "App", results: []*api.MetaRecord{records[0].record, records[1].record, records[4].record}},
		{term: "hello", results: []*api.MetaRecord{
			records[2].record, records[5].record,
		}},
		{term: "text", results: []*api.MetaRecord{
			records[3].record, records[5].record,
		}},
		{term: "Indexer", results: []*api.MetaRecord{records[5].record}},
		{term: "indexer", results: []*api.MetaRecord{records[5].record}},
	}

	for _, d := range data {
		results, _ := index.Search(d.term)
		require.ElementsMatchf(t, d.results, results, "index should return the correct results for the following query: %s", d.term)
	}
}

func TestExactMatchSearching(t *testing.T) {
	index, err := newIndex(false)
	require.NoError(t, err)

	for _, r := range records {
		index.Index(r.record, r.data)
	}

	_, err = index.Search("")
	require.Error(t, err, "search should fail if provided an empty term")

	data := []struct {
		term    string
		results []*api.MetaRecord
	}{
		{term: "App 1", results: []*api.MetaRecord{records[0].record}},
		{term: "App 2", results: []*api.MetaRecord{records[1].record}},
		{term: "hello@email.com", results: []*api.MetaRecord{records[2].record}},
		{term: "This is some sample text", results: []*api.MetaRecord{records[3].record}},
	}

	for _, d := range data {
		results, _ := index.Search(d.term)
		require.ElementsMatchf(t, d.results, results, "index should return the correct results for the following query: %s", d.term)
	}
}
