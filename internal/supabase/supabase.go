package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/joyfuldevs/project-jarvis/internal/setting"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type Client struct {
	Domain  string
	APIKey  string
	AuthKey string

	HTTPClient *http.Client
}

// GetJarvisSetting retrieves the Jarvis settings for a given channel ID.
func (c *Client) GetJarvisSetting(ctx context.Context, channelID string) (*setting.JarvisSetting, error) {
	data, err := c.doRequest(ctx, apiRequest{
		method: "GET",
		path:   path.Join("/rest/v1/jarvis_settings"),
		params: map[string]string{"id": "eq." + channelID},
		body:   nil,
	})
	if err != nil {
		return nil, err
	}

	var settings []*setting.JarvisSetting
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	if len(settings) == 0 {
		return nil, ErrNotFound
	}

	return settings[0], nil
}

func (c *Client) SetJarvisSetting(ctx context.Context, s *setting.JarvisSetting) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = c.doRequest(ctx, apiRequest{
		method: "POST",
		path:   path.Join("/rest/v1/jarvis_settings"),
		body:   bytes.NewReader(data),
	})

	return err
}

type apiRequest struct {
	method string
	path   string
	params map[string]string
	body   io.Reader
}

func (c *Client) doRequest(ctx context.Context, r apiRequest) ([]byte, error) {
	apiURL := url.URL{Scheme: "https", Host: c.Domain, Path: r.path}
	query := url.Values{}
	for k, v := range r.params {
		query.Set(k, v)
	}
	apiURL.RawQuery = query.Encode()
	fmt.Println(apiURL.String())
	req, err := http.NewRequestWithContext(ctx, r.method, apiURL.String(), r.body)
	if err != nil {
		return nil, err
	}

	// API key and Auth key headers
	req.Header.Add("apikey", c.APIKey)
	req.Header.Add("Authorization", "Bearer "+c.AuthKey)
	req.Header.Add("Prefer", "resolution=merge-duplicates")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &url.Error{Op: r.method, URL: apiURL.String(), Err: errors.New(resp.Status)}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
