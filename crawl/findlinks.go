package crawl

import (
	"fmt"
	"regexp"
)

func CountLinks(html string) (int, error) {
	matches, err := c.FindLinks(html)
	if err != nil {
		fmt.Errorf("failed to find links: %w", err)
	}
	return len(matches), nil
}

func (c *Crawler) FindLinks(html string) ([]string, error) {
	anchorTagRegex := `<a\s+[^>]*href=[^>]*>`
	re, err := regexp.Compile(anchorTagRegex)
	if err != nil {
		return nil, fmt.Errorf("Failed to compile re: %w", err)
	}
	matches := re.FindAll([]byte(html), -1)
	links := string(matches)
	for _, match := range matches {
		match = []byte(c.resolveLink(string(match)))
	}
	return matches, nil
}
