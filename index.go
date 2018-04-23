package main

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/olivere/elastic"
)

const (
	manpageIndexName   = "elastic_manpages"
	aproposRegex       = `(.*)\(.*\s*-\s*(.*)`
	concurrentIndexers = 30
)

type Manpage struct {
	Command     string `json:"command,omitempty"`
	Description string `json:"description,omitempty"`
	Text        string `json:"manpage,omitempty"`
}

func createIndex(client *elastic.Client) {
	removeExistingManpageIndex(client)
	createManpageIndex(client)
	pages := getManpages()
	indexEachPage(pages, client)
}

func indexEachPage(pages []string, client *elastic.Client) {
	log.Print(len(pages))

	rateLimiter := make(chan bool, concurrentIndexers)
	var wg sync.WaitGroup
	wg.Add(len(pages))

	for number, page := range pages {
		go func(page string, number int) {
			rateLimiter <- true
			defer func() { <-rateLimiter }()

			log.Println(number)
			newPages := handleMultipleCommandsOnSameLine(page)

			for _, newPage := range newPages {

				re := regexp.MustCompile(aproposRegex)
				submatches := re.FindStringSubmatch(newPage)

				if len(submatches) < 3 {
					log.Printf("Error finding matches for %s", newPage)
					continue
				}

				command := submatches[1]
				description := submatches[2]
				manpage := strconv.Quote(getManpage(submatches[1]))
				manpageStruct := Manpage{
					Command:     "\"" + command + "\"",
					Description: "\"" + description + "\"",
					Text:        manpage,
				}

				_, err := client.Index().
					Index(manpageIndexName).
					Type("document").
					BodyJson(manpageStruct).
					Do(context.TODO())

				if err != nil {
					log.Printf("Error indexing page: %s, %s", page, err)
				}
			}

			wg.Done()
		}(page, number)
	}

	wg.Wait()
	log.Println("Indexing done!!!")

}

func handleMultipleCommandsOnSameLine(line string) []string {
	hyphenIndex := strings.Index(line, "-")

	if hyphenIndex == -1 {
		log.Printf("No hyphen in line: %s", line)
		return []string{}
	}

	description := line[hyphenIndex:]

	splitOnComma := strings.Split(line[:hyphenIndex], ",")
	for i, splitted := range splitOnComma {
		splitOnComma[i] = strings.TrimSpace(splitted) + " " + description
	}

	return splitOnComma
}

func getManpages() []string {
	cmd := exec.Command("apropos", ".")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	outString := out.String()
	return strings.Split(outString, "\n")
}

func getManpage(name string) string {
	cmd := exec.Command("man", name)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		// Try to run the command again but with a lowercase manpage name.
		cmd = exec.Command("man", strings.ToLower(name))
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			log.Printf("Error getting manpage: %s, %s", name, err)
			return ""
		}
	}

	return out.String()
}

func removeExistingManpageIndex(client *elastic.Client) {
	exists, err := client.IndexExists(manpageIndexName).Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		_, err = client.DeleteIndex(manpageIndexName).Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createManpageIndex(client *elastic.Client) {
	_, err := client.CreateIndex(manpageIndexName).
		Body(
			`{"mappings": {
				"document": {
					"properties": {
						"command": {
							"type": "text"
						},
						"description": {
							"type": "text",
							"analyzer": "english"
						},
						"manpage": {
							"type": "text",
							"analyzer": "english"
						}
					}
				}
			}}`).
		Do(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
}
