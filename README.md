# manpage-elastic [![Build Status](https://travis-ci.org/ephraimkunz/manpage-elastic.svg?branch=master)](https://travis-ci.org/ephraimkunz/manpage-elastic)
Inspired by [iridakos's tutorial](https://iridakos.com/tutorials/2018/04/12/elasticsearch-linux-man-pages.html), but ported to Go with additional features.

## Usage
1. Make sure Elasticsearch cluster is running on `http://localhost:9200`.
2. Run webserver UI with `go run main.go`.

## Additional Features
* Concurrent index creation with configurable number of goroutines
* Command-line flags to force creation of an index and do health checks
* API and basic web frontend
* Comes with unit tests
