package catapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	http    *http.Client
	baseURL string
	apiKey  string
}
type Breed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewClient(baseURL, apiKey string, timeout time.Duration) *Client {
	return &Client{
		http:    &http.Client{Timeout: timeout},
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
	}
}
func (c *Client) SearchBreeds(ctx context.Context, q string) ([]Breed, int, error) {
	if strings.TrimSpace(q) == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("empty query")
	}

	endpoint := fmt.Sprintf("%s/v1/breeds/search?q=%s", c.baseURL, url.QueryEscape(q))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}
 

	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("catapi responded with status %d", resp.StatusCode)
	}

	var breeds []Breed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("decode response: %w", err)
	}

	return breeds, resp.StatusCode, nil
}
