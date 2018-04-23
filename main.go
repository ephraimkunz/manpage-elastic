package main

import (
	"context"
	"flag"
	"log"

	"github.com/olivere/elastic"
)

var shouldIndex = flag.Bool("index", false, "Should man pages be indexed into ElasticSearch?")
var searchText = flag.String("query", "", "Query to search for in man pages")
var verbose = flag.Bool("verbose", false, "Print verbose output (cluster status)")

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
		createIndex(client)
	}

	runSearch(client, *searchText)
}
