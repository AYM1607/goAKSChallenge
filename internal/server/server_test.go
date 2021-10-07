package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/AYM1607/goAKSChallenge/api"
	"github.com/AYM1607/goAKSChallenge/internal/common"
	"github.com/AYM1607/goAKSChallenge/internal/server"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const (
	// Invalid files.
	invRecordsDir = "testdata/invalid"

	// Valid files.
	// The file names reflect the content and can be used as reference to know
	// what terms to use in search tests.
	record1Fp = "testdata/valid/ValApp1Maint1Web1Desc2.yaml"
	record2Fp = "testdata/valid/ValApp2Maint1Web1Desc1.yaml"
	record3Fp = "testdata/valid/ValApp3Maint2Web2Desc1.yaml"
	record4Fp = "testdata/valid/ValApp4Maint2Web2Desc3.yaml"
)

var recordsFps = []string{record1Fp, record2Fp, record3Fp, record4Fp}

func createServer(t *testing.T) *httptest.Server {
	h, err := server.NewHTTPHandler()
	require.NoError(t, err, "test server should be able to be created correclty")

	server := httptest.NewServer(h)
	// Avoid deferring server cleanup in main test function.
	t.Cleanup(func() {
		server.Close()
	})
	return server
}

func TestInvalidCreates(t *testing.T) {
	filePaths, err := common.GetAllFilesInDir(invRecordsDir)
	require.NoError(t, err, "invalid test data directory should be present")

	server := createServer(t)
	e := httpexpect.New(t, server.URL)

	for _, fp := range filePaths {
		fb, err := os.ReadFile(fp)
		require.NoError(t, err, "testdata file should be able to be opened successfully")

		reqBody := map[string]string{
			"record": string(fb),
		}
		e.POST("/records").WithJSON(reqBody).
			Expect().
			Status(http.StatusBadRequest)
	}
}

func TestValidCreatesAndSearches(t *testing.T) {
	records := map[string]string{}

	// Open the test files and parse them to strings to be used in requests and comparisons.
	for _, rFp := range recordsFps {
		rb, err := os.ReadFile(rFp)
		require.NoError(t, err, "testdata file should be able to be opened successfully.")
		records[rFp] = string(rb)
	}

	testServer := createServer(t)
	e := httpexpect.New(t, testServer.URL)

	// Test successful creation fo valid records.
	for k, r := range records {
		t.Run(fmt.Sprintf("Create/%s", k), func(_ *testing.T) {
			reqBody := map[string]string{
				"record": r,
			}
			e.POST("/records").WithJSON(reqBody).
				Expect().
				Status(http.StatusCreated)
		})
	}

	searchTestsData := []struct {
		name    string
		req     server.SearchRequest
		results []string
	}{
		{
			name:    "Exact, one result",
			req:     server.SearchRequest{JoinMethod: "or", SearchTerms: []api.SearchTerm{{Field: "title", Query: "Valid App 1"}}},
			results: []string{records[record1Fp]},
		},
		{
			name: "Exact, multiple results, OR join",
			req: server.SearchRequest{JoinMethod: "or", SearchTerms: []api.SearchTerm{
				{Field: "title", Query: "Valid App 1"},
				{Field: "maintainerEmail", Query: "man2@mail.com"},
			}},
			results: []string{records[record1Fp], records[record3Fp], records[record4Fp]},
		},
		{
			name: "Exact, multiple results, AND join",
			req: server.SearchRequest{JoinMethod: "and", SearchTerms: []api.SearchTerm{
				{Field: "website", Query: "https://website2.io"},
				{Field: "maintainerEmail", Query: "man2@mail.com"},
			}},
			results: []string{records[record3Fp], records[record4Fp]},
		},
		{
			name: "FTS, one result",
			req: server.SearchRequest{JoinMethod: "or", SearchTerms: []api.SearchTerm{
				{Field: "description", Query: "threeForTesting"},
			}},
			results: []string{records[record4Fp]},
		},
		{
			name: "FTS, multiple results, OR join",
			req: server.SearchRequest{JoinMethod: "or", SearchTerms: []api.SearchTerm{
				{Field: "description", Query: "twoForTesting"},
				{Field: "description", Query: "threeForTesting"},
			}},
			results: []string{records[record1Fp], records[record4Fp]},
		},
		{
			name: "FTS, multiple results, AND join",
			req: server.SearchRequest{JoinMethod: "and", SearchTerms: []api.SearchTerm{
				{Field: "description", Query: "description"},
				{Field: "description", Query: "oneForTesting"},
			}},
			results: []string{records[record2Fp], records[record3Fp]},
		},
		{
			name: "FTS and exact, one result",
			req: server.SearchRequest{JoinMethod: "and", SearchTerms: []api.SearchTerm{
				{Field: "description", Query: "threeForTesting"},
				{Field: "title", Query: "Valid App 4"},
			}},
			results: []string{records[record4Fp]},
		},
		{
			name: "FTS and exact, multiple results, OR join",
			req: server.SearchRequest{JoinMethod: "or", SearchTerms: []api.SearchTerm{
				{Field: "description", Query: "oneForTesting"},
				{Field: "title", Query: "Valid App 4"},
			}},
			results: []string{records[record2Fp], records[record3Fp], records[record4Fp]},
		},
		{
			name: "FTS and exact, multiple results, AND join",
			req: server.SearchRequest{JoinMethod: "and", SearchTerms: []api.SearchTerm{
				{Field: "website", Query: "https://website1.io"},
				{Field: "description", Query: "description"},
			}},
			results: []string{records[record1Fp], records[record2Fp]},
		},
	}

	for _, td := range searchTestsData {
		t.Run(td.name, func(t *testing.T) {
			// httpexpect is built in a functional-is way so the call to Status returns the original response.
			rawRes := e.POST("/records/search").WithJSON(td.req).
				Expect().
				Status(http.StatusOK)

			rawBody := rawRes.Body().Raw()
			res := server.SearchResponse{}

			err := json.Unmarshal([]byte(rawBody), &res)
			require.NoError(t, err, "successfull search requests should be unmarshable")
			require.ElementsMatch(t, res.Records, td.results, "the api should return the expected results.")
		})
	}
}
