package crawl

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

func NewCrawler(baseURL string, client *http.Client) (*Crawler, error) {
	URL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %q url: %w", baseURL, err)
	}
	return &Crawler{
		URL,
		client,
	}, nil
}

type WebsiteInfo struct {
	Html      string
	LinkCount int
}

func (c *Crawler) resolveLink(link string) (string, error) {
	parsedLink, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("invalid link: %s, error: %w", link, err)
	}
	return c.baseURL.ResolveReference(parsedLink).String(), nil
}

func (c *Crawler) Crawl() <-chan WebsiteInfo {
	websiteStream := make(chan WebsiteInfo)
	linkStream := make(chan string)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer close(websiteStream)
		c.crawlSite(&wg, websiteStream, linkStream, c.baseURL.String())
		wg.Wait()
		close(linkStream)
	}()

	go func() {
		select {
		case link, ok := <-linkStream:
			if !ok {
				return
			}
			wg.Add(1)
			c.crawlSite(&wg, websiteStream, linkStream, link)
		}
	}()

	return websiteStream
}

func (c *Crawler) crawlSite(wg *sync.WaitGroup, websiteStream chan WebsiteInfo, linkStream chan string, url string) {
	defer wg.Done()
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}
	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("Failed to Do request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		return
	}

	html := string(body)

	links, err := c.FindLinks(html)
	if err != nil {
		log.Printf("failed to find links: %v", err)
		return
	}
	linkCount := len(links)

	for _, link := range links {
		linkStream <- string(link)
	}

	websiteStream <- WebsiteInfo{
		html,
		linkCount,
	}
}
