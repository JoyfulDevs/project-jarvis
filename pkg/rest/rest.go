package rest

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Domain string
}

func NewClient(domain string) *Client {
	return &Client{Domain: domain}
}

func (c *Client) RequestAPI(
	ctx context.Context,
	method string,
	path string,
	opts ...Option,
) ([]byte, error) {
	options := &options{
		httpClient: http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(options)
	}

	url := c.Domain + path + options.params
	req, err := http.NewRequestWithContext(ctx, method, url, options.body)
	if err != nil {
		return nil, err
	}

	for key, value := range options.headers {
		req.Header.Add(key, value)
	}

	resp, err := options.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if len(data) > 0 {
			return nil, errors.New(string(data))
		} else {
			return nil, errors.New(resp.Status)
		}
	}

	return data, nil
}
