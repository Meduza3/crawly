package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/meduza3/crawly/crawl"
)

func main() {
	client := http.Client{}
	crawler, err := crawl.NewCrawler("https://cs.pwr.edu.pl/gebala/dyd/", &client)
	if err != nil {
		log.Fatalf("error creating crawler: %v", err)
	}
	websiteInfo := crawler.Crawl()
	for website := range websiteInfo {
		fmt.Println(website)
	}
}
