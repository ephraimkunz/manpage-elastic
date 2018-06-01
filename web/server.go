package web

import (
	"html/template"
	"net/http"

	"github.com/ephraimkunz/manpage-elastic/search"
	"github.com/olivere/elastic"
)

type SearchHandler struct {
	Command         func(query string) (*search.SearchResults, error)
	ResultsTemplate *template.Template
	WelcomeTemplate *template.Template
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

func NewSearchHandler(command func(string) (*search.SearchResults, error)) SearchHandler {
	handler := SearchHandler{}
	handler.Command = command
	handler.ResultsTemplate = template.Must(template.ParseFiles("web/templates/results.html"))
	handler.WelcomeTemplate = template.Must(template.ParseFiles("web/templates/welcome.html"))
	return handler
}

func (handler *SearchHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	handler.WelcomeTemplate.Execute(w, nil)
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

	handler.ResultsTemplate.Execute(w, sr)
}
