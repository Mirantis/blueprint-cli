package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// ReadURI reads the content of a URI.
// The URI argument can be a file path or a URL.
// @TODO: Make this function testable by injecting a reader for file and http requests.
func ReadURI(uri string) ([]byte, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URI: %w", err)
	}

	switch u.Scheme {
	case "http", "https":
		return readFromUrl(uri)
	default:
		filePath := strings.Replace(uri, "file://", "", 1)
		return readFromPath(filePath)
	}
}

func readFromPath(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return content, nil
}

func readFromUrl(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read URL body: %w", err)
	}
	return body, nil
}
