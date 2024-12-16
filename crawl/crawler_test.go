package crawl_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/meduza3/crawly/crawler"
)

func TestCrawler(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		expectedHTML   string
		expectError    bool
	}{
		{
			name:           "Successful crawl with valid HTML",
			serverResponse: "<html><body>Hello, World!</body></html>",
			expectedHTML:   "<html><body>Hello, World!</body></html>",
			expectError:    false,
		},
		{
			name:           "Empty response",
			serverResponse: "",
			expectedHTML:   "",
			expectError:    false,
		},
		{
			name:           "Invalid HTML",
			serverResponse: "<html><body><p>Missing closing tags",
			expectedHTML:   "<html><body><p>Missing closing tags",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable????
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(tt.serverResponse))
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer mockServer.Close()

			parsedURL, err := url.Parse(mockServer.URL)
			if err != nil {
				t.Fatalf("Failed to parse mock server URL: %v", err)
			}

			// Initialize the Crawler with the mock server's URL and HTTP client.
			crawler := crawler.NewCrawler(parsedURL, mockServer.Client())

			// Call the Crawl method.
			websiteChan := crawler.Crawl()

			// Receive from the channel.
			websiteInfo, ok := <-websiteChan
			if !ok {
				t.Fatalf("Expected to receive WebsiteInfo, but channel was closed")
			}

			// Verify the HTML content.
			if websiteInfo.Html != tt.expectedHTML {
				t.Errorf("Expected HTML: %q, got: %q", tt.expectedHTML, websiteInfo.Html)
			}
		})
	}
}

func TestCountLinks(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected int
		hasError bool
	}{
		{
			name:     "Valid HTML with multiple links",
			html:     `<html><a href="https://example.com">Link</a><a href="https://example.org">Another</a></html>`,
			expected: 2,
			hasError: false,
		},
		{
			name:     "HTML with no links",
			html:     `<html><p>No links here</p></html>`,
			expected: 0,
			hasError: false,
		},
		{
			name:     "Malformed HTML with one valid link",
			html:     `<html><a href="https://example.com">Link<p><a href="https://example.org"></html>`,
			expected: 2,
			hasError: false,
		},
		{
			name:     "Empty HTML string",
			html:     ``,
			expected: 0,
			hasError: false,
		},
		{
			name:     "Valid HTML with links missing href",
			html:     `<html><a>Missing href</a><a name="anchor">Anchor</a></html>`,
			expected: 0,
			hasError: false,
		},
		{
			name:     "HTML with nested links",
			html:     `<html><a href="https://example.com"><a href="https://nested.com"></a></a></html>`,
			expected: 2,
			hasError: false,
		},
		{
			name:     "HTML with links having extra attributes",
			html:     `<html><a href="https://example.com" class="link">Example</a></html>`,
			expected: 1,
			hasError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			count, err := crawler.CountLinks(test.html)
			if (err != nil) != test.hasError {
				t.Errorf("Unexpected error status. Got: %v, Expected error: %v", err != nil, test.hasError)
			}
			if count != test.expected {
				t.Errorf("Unexpected link count. Got: %d, Expected: %d", count, test.expected)
			}
		})
	}
}
