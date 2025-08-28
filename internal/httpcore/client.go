package httpcore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseUrl string
	Client  *http.Client
	Headers map[string]string
}

// New creates a new HTTP client with the specified base URL, timeout, and headers.
func New(baseUrl string, timeout time.Duration, headers map[string]string) *Client {
	return &Client{
		BaseUrl: strings.TrimSuffix(baseUrl, "/"),
		Client: &http.Client{
			Timeout: timeout,
		},
		Headers: headers,
	}
}

// Get sends a GET request to the specified URL and decodes the JSON response into the provided output variable.
func (c *Client) Get(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+url, nil)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *Client) do(req *http.Request, out any) error {
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s %s: %s - %s", req.Method, req.URL.Path, resp.Status, string(b))
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}
