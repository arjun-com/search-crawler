package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"example.com/web_crawler/utils"
)

var visited = make(map[string]bool)
var contained []string

var visitedMutex = sync.Mutex{}
var waitGroup = sync.WaitGroup{}
var sem = make(chan struct{}, 10)

func crawl(startUrl string, phrase string) {
	defer waitGroup.Done()

	sem <- struct{}{}
	defer func() {
		<-sem
	}()

	if len(contained)%100 == 0 {
		log.Println("Searched ", len(contained), " webpages.")
	}

	resp, err := http.Get(startUrl)
	if err != nil {
		log.Println("ERROR in GET request.\n", err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("ERROR reading body.\n", err.Error())
		return
	}

	if strings.Contains(string(body), phrase) {
		log.Println(startUrl, " contains search phrase.")
		contained = append(contained, startUrl)
	}

	links, err := utils.GetLinks(string(body))
	if err != nil {
		log.Fatal("ERROR parsing links.\n", err.Error())
		return
	}

	for _, link := range links {
		url, err := utils.ResolveUri(startUrl, link)
		if err != nil {
			log.Println("ERROR resolving BASEURL: ", startUrl, " with URI: ", link)
			continue
		}

		visitedMutex.Lock()
		_, seen := visited[url]
		if !seen {
			visited[url] = true
			visitedMutex.Unlock()
			waitGroup.Add(1)
			go crawl(url, phrase)
			continue
		}
		visitedMutex.Unlock()
	}
}

func main() {
	startUrl := flag.String("u", "", "The starting URL to crawl")
	phrase := flag.String("p", "", "The phrase to search for in pages")

	flag.Parse()

	if *startUrl == "" || *phrase == "" {
		fmt.Println("Usage: go run main.go -u=http://example.com -p=example")
		return
	}

	waitGroup.Add(1)
	go crawl(*startUrl, *phrase)
	waitGroup.Wait()

	log.Println("Done Crawling\n")

	for _, link := range contained {
		fmt.Println(link)
	}
}
