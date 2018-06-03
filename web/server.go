package web

import (
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ephraimkunz/manpage-elastic/search"
	"github.com/olivere/elastic"
)

const (
	defaultNumResultsToReturn = 10
)

type SearchHandler struct {
	Command         func(query string, count int64) (*search.SearchResults, error)
	ResultsTemplate *template.Template
	WelcomeTemplate *template.Template
}

type SearchCreator struct {
}

func (sc *SearchCreator) Run(query string, count int64) (*search.SearchResults, error) {
	client, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}

	return search.RunSearch(client, query, count)
}

func NewSearchHandler(command func(string, int64) (*search.SearchResults, error)) SearchHandler {
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
	count, err := countFromQuery(r.URL.Query())
	if err != nil {
		http.Error(w, "Bad query string", http.StatusBadRequest)
		return
	}

	sr, err := handler.Command(query, count)
	if err != nil {
		// render an error
		http.Error(w, "Failed to get search results", http.StatusInternalServerError)
		return
	}

	handler.ResultsTemplate.Execute(w, sr)
}

func countFromQuery(query url.Values) (int64, error) {
	countString := query.Get("count")

	if countString == "" {
		return defaultNumResultsToReturn, nil
	}

	count, err := strconv.ParseInt(countString, 10, 32)
	return count, err
}
