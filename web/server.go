package web

import (
	"encoding/json"
	"net/http"

	"github.com/ephraimkunz/manpage-elastic/search"
	"github.com/olivere/elastic"
)

type SearchHandler struct {
	Command func(query string) (*search.SearchResults, error)
}

type SearchCreator struct {
}

func (sc *SearchCreator) Run(query string) (*search.SearchResults, error) {
	client, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}

	return search.RunSearch(client, query)
}

func (handler *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	// 1. parse user data
	query := r.URL.Query().Get("query")
	sr, err := handler.Command(query)
	if err != nil {
		// render an error
		http.Error(w, "Failed to get search results", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sr)
	return

}
