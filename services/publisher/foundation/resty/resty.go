package resty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

type Payload struct {
	data        []byte
	contentType string
}

func NewJSONPayload(v interface{}) (*Payload, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON payload: %w", err)
	}
	return &Payload{data: data, contentType: "application/json"}, nil
}

func (p *Payload) Reader() io.Reader {
	return bytes.NewReader(p.data)
}

func (p *Payload) ContentType() string {
	return p.contentType
}

func (p *Payload) String() string {
	return string(p.data)
}

func New(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{Timeout: timeout},
	}
}

func (c *Client) Request(ctx context.Context, method string, path string, headers http.Header, payload *Payload, result interface{}) error {
	url := c.baseURL + path

	var body io.Reader
	if payload != nil {
		body = payload.Reader()
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	if headers != nil {
		req.Header = headers
	}

	if payload != nil {
		req.Header.Set("Content-Type", payload.ContentType())
	}

	log.Printf("Making %s request to %s", method, url)
	log.Printf("Headers: %+v", req.Header)
	if payload != nil {
		log.Printf("Payload: %s", payload.String())
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Response status: %s", resp.Status)

	if err := handleResponse(resp, result); err != nil {
		return fmt.Errorf("error handling response: %w", err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, path string, headers http.Header, result interface{}) error {
	return c.Request(ctx, http.MethodGet, path, headers, nil, result)
}

func (c *Client) Post(ctx context.Context, path string, headers http.Header, payload *Payload, result interface{}) error {
	return c.Request(ctx, http.MethodPost, path, headers, payload, result)
}

func handleResponse(resp *http.Response, result interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	log.Printf("Response body: %s", string(body))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	if result == nil {
		return nil
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("error unmarshaling response: %w", err)
	}

	return nil
}
