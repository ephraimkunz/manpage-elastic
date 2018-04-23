# manpage-elastic
Inspired by [iridakos's tutorial](https://iridakos.com/tutorials/2018/04/12/elasticsearch-linux-man-pages.html), but ported to Go with additional features.

## Usage
1. Build project with `go build`. `dep` is used for vendoring.
2. Make sure ElasticSearch cluster is running on `http://localhost:9200`.
3. ./manpage-elastic [-verbose] [-index] -query \<query string\>

## Additional Features
* Concurrent indexes with configurable number of goroutines
* Commandline flags to force creation of an index and do health checks
