# manpage-elastic [![Build Status](https://travis-ci.org/ephraimkunz/manpage-elastic.svg?branch=master)](https://travis-ci.org/ephraimkunz/manpage-elastic)
Inspired by [iridakos's tutorial](https://iridakos.com/tutorials/2018/04/12/elasticsearch-linux-man-pages.html), but ported to Go with additional features.

## Usage
1. Build project with `go build`. [dep](https://github.com/golang/dep) is used for vendoring.
2. Make sure Elasticsearch cluster is running on `http://localhost:9200`.
3. ./manpage-elastic [-verbose] [-index] -query \<query string\>

## Additional Features
* Concurrent index creation with configurable number of goroutines
* Command-line flags to force creation of an index and do health checks
