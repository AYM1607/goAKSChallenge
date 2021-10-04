package api

type searchField uint8

type SearchTerm struct {
	Field string `json:"field"`
	Query string `json:"query"`
}
