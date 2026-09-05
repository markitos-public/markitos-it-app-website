package faqs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type FAQ struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (client *Client) List() ([]FAQ, error) {
	request, err := http.NewRequest(http.MethodGet, client.baseURL+"/faqs", nil)
	if err != nil {
		return nil, fmt.Errorf("create FAQ request: %w", err)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request FAQs: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		return []FAQ{}, nil
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FAQ API returned status %s", response.Status)
	}

	result := make([]FAQ, 0)
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		if err == io.EOF {
			return result, nil
		}
		return nil, fmt.Errorf("decode FAQs response: %w", err)
	}

	return result, nil
}
