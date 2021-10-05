package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AYM1607/goAKSChallenge/api"
	"github.com/AYM1607/goAKSChallenge/internal/store"
	"github.com/goccy/go-yaml"
	"github.com/gorilla/mux"
)

const searchErrString = "search request could not be completed due to an internal error"

func New(addr string) *http.Server {
	handler := newHandler()

	r := mux.NewRouter()

	r.HandleFunc("/records", handler.handleCreate).Methods("POST")
	// We could debate using POST or GET for a search endpoint. For this challenge I'll prioritize ease of parsing.
	// Since the GET verb does not support a body, we would need to parse search terms from the URL.
	// If the requirements mentioned compatibility with browsers or ease of query sharing the effort of using
	// query string params would be justified.
	r.HandleFunc("/records/search", handler.handleSearch).Methods("POST")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type handler struct {
	Store *store.Store
}

func newHandler() *handler {
	return &handler{
		Store: store.New(),
	}
}

type CreateRequest struct {
	Record string `json:"record"`
}

type CreateResponse struct {
	Message string `json:"message"`
}

type SearchRequest struct {
	JoinMethod  api.SearchJoinMethod `json:"joinMethod"`
	SearchTerms []api.SearchTerm     `json:"searchTerms"`
}

// Since the yaml is accepted as a string, the records that are found from a search
// are also returned as strings.
type SearchResponse struct {
	Records []string `json:"records"`
}

func (h *handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure the payload is valid.
	err = h.Store.Append([]byte(req.Record))
	if err != nil {
		// This error string contains information about what went wrong with the payload processing,
		// including field names that caused the error.
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := CreateResponse{Message: "The record was added successfully."}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		// Generic error, since the response is always the same and is always
		// encodable, this should not happen but leaving it as a safeguard.
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (h *handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	var req SearchRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure payload is valid.
	if err := req.JoinMethod.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// NOTE: This validation is linear in terms of time complexity,
	// beware of search requests with a high number of terms.
	invalidFields := []string{}
	for _, term := range req.SearchTerms {
		if err := term.Field.IsValid(); err != nil {
			invalidFields = append(invalidFields, string(term.Field))
		}
	}

	if len(invalidFields) > 0 {
		http.Error(w,
			fmt.Sprintf("the following field(s) are not supported: %s", strings.Join(invalidFields, ",")),
			http.StatusBadRequest)
		return
	}

	records, err := h.Store.Search(req.JoinMethod, req.SearchTerms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawRecords := []string{}
	for _, record := range records {

		rawRecord, err := yaml.Marshal(record)
		// Since all records where unmarshalled from valid yaml this should not
		// happen but leaving it as a safeguard.
		if err != nil {
			http.Error(w, searchErrString, http.StatusInternalServerError)
			return
		}
		rawRecords = append(rawRecords, string(rawRecord))
	}

	res := SearchResponse{Records: rawRecords}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&res)
	if err != nil {
		http.Error(w, searchErrString, http.StatusInternalServerError)
	}
}
