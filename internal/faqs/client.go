package faqs

import (
	"encoding/json"
	"fmt"
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
	endpoint   string
	httpClient *http.Client
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: strings.TrimRight(endpoint, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (client *Client) List() ([]FAQ, error) {
	request, err := http.NewRequest(http.MethodGet, client.endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create FAQ request: %w", err)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request FAQs: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FAQ API returned status %s", response.Status)
	}

	var result []FAQ
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode FAQs response: %w", err)
	}

	return result, nil
}
