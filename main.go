package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/ephraimkunz/manpage-elastic/search"
	"github.com/ephraimkunz/manpage-elastic/web"
	"github.com/olivere/elastic"
)

var shouldIndex = flag.Bool("index", false, "Should man pages be indexed into ElasticSearch?")
var searchText = flag.String("query", "", "Query to search for in man pages")
var verbose = flag.Bool("verbose", false, "Print verbose output (cluster status)")
var runServer = flag.Bool("runServer", true, "Runs a server for remote queries")

func main() {
	flag.Parse()

	client, err := elastic.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	if *verbose {
		health, err := client.ClusterHealth().Pretty(true).Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Cluster status: %s", health.Status)
	}

	if *shouldIndex {
		search.CreateIndex(client)
	}

	if len(*searchText) > 0 {
		search.RunSearch(client, *searchText)
	}

	if *runServer {
		searchCreator := web.SearchCreator{}
		searchHandler := web.NewSearchHandler(searchCreator.Run)
		http.HandleFunc("/search", searchHandler.Search)
		http.HandleFunc("/", searchHandler.Welcome)

		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}
