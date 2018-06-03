package search

import (
	"context"
	"reflect"
	"strings"

	"github.com/olivere/elastic"
)

type SearchResults struct {
	NumHits int64 `json:"num_hits,omitempty"`
	Hits    []Hit `json:"hits,omitempty"`
}

type Hit struct {
	Command     string `json:"command,omitempty"`
	Description string `json:"description,omitempty"`
	Ordinal     int    `json:"ordinal,omitempty"`
}

func RunSearch(client *elastic.Client, search string, count int64) (*SearchResults, error) {
	q := elastic.NewMultiMatchQuery(search, "command", "description^3", "manpage^3").
		Operator("or").
		TieBreaker(1.0).
		CutoffFrequency(0.1).
		Type("cross_fields")

	res, err := client.Search(manpageIndexName).
		Size(int(count)).
		Query(q).
		Do(context.TODO())

	if err != nil {
		return nil, err
	}

	return createSearchResults(res, search), nil
}

func createSearchResults(res *elastic.SearchResult, search string) *SearchResults {
	var ttyp Manpage
	results := &SearchResults{NumHits: res.TotalHits()}

	for idx, item := range res.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Manpage); ok {
			hit := Hit{Ordinal: idx + 1, Command: strings.Trim(t.Command, `"`), Description: strings.Trim(t.Description, `"`)}
			results.Hits = append(results.Hits, hit)
		}
	}

	return results
}
