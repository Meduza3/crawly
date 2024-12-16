package crawler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type Crawler struct {
	baseURL *url.URL
	client  *http.Client
}

func NewCrawler(baseURL *url.URL, client *http.Client) *Crawler {
	return &Crawler{
		baseURL,
		client,
	}
}

type WebsiteInfo struct {
	Html string
}

func (c *Crawler) Crawl() <-chan WebsiteInfo {
	websiteStream := make(chan WebsiteInfo)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		c.crawlSite(&wg, websiteStream, c.baseURL.String())
	}()

	wg.Wait()
	return websiteStream
}

func (c *Crawler) crawlSite(wg *sync.WaitGroup, websiteStream chan WebsiteInfo, url string) {
	defer wg.Done()
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("Failed to Do request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	html := string(body)
	fmt.Println(html)
	websiteStream <- WebsiteInfo{
		html,
	}
}
