package crawler_test

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
