package server

import (
	"encoding/json"
	"net/http"

	"github.com/AYM1607/goAKSChallenge/api"
	"github.com/AYM1607/goAKSChallenge/internal/store"
	"github.com/gorilla/mux"
)

func New(addr string) *http.Server {
	handler := newHandler()

	r := mux.NewRouter()

	r.HandleFunc("/records", handler.handleCreate).Methods("POST")

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
	SearchTerms []api.SearchTerm `json:"searchTerms"`
}

type SearchResponse struct {
	Records []string `json:"records"`
}

type DeepSearchRequest api.SearchTerm

func (h *handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	// Ensure the request is formed correctly.
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

}
