package main

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/olivere/elastic"
)

func runSearch(client *elastic.Client, search string) {
	q := elastic.NewMultiMatchQuery(search, "command", "description^3", "manpage^3").
		Operator("or").
		TieBreaker(1.0).
		CutoffFrequency(0.1).
		Type("cross_fields")

	res, err := client.Search(manpageIndexName).
		Size(10).
		Query(q).
		Do(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	prettyPrintSearchResults(res, search)
}

func prettyPrintSearchResults(res *elastic.SearchResult, search string) {
	var ttyp Manpage

	// Get max widths of command and description
	maxCommand, maxDescription := 0, 0
	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Manpage); ok {
			if len(t.Command) > maxCommand {
				maxCommand = len(t.Command)
			}

			if len(t.Description) > maxDescription {
				maxDescription = len(t.Description)
			}
		}
	}

	formatString := fmt.Sprintf("%%-%ds: %%-%ds\n", maxCommand+1, maxDescription)

	fmt.Printf(`--- Results matching query "%s"---`, search)
	fmt.Println()
	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Manpage); ok {
			fmt.Printf(formatString, t.Command, t.Description)
		}
	}
}
