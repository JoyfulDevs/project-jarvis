package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/joyfuldevs/project-jarvis/internal/setting"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrInvalidSetting = errors.New("invalid jarvis setting")
)

type Client struct {
	Domain  string
	APIKey  string
	AuthKey string

	HTTPClient *http.Client
}

// GetJarvisSetting 는 주어진 채널 ID에 대한 Jarvis 설정을 가져옵니다.
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

// SetJarvisSetting 는 주어진 Jarvis 설정을 저장 하거나 갱신합니다.
func (c *Client) SetJarvisSetting(ctx context.Context, s *setting.JarvisSetting) error {
	if s == nil || s.ID == "" {
		return ErrInvalidSetting
	}

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
	req, err := http.NewRequestWithContext(ctx, r.method, apiURL.String(), r.body)
	if err != nil {
		return nil, err
	}

	// API 키와 인증 키를 헤더에 추가
	req.Header.Add("apikey", c.APIKey)
	req.Header.Add("Authorization", "Bearer "+c.AuthKey)
	// 이미 데이터가 있는 경우 UPSERT로 동작 하도록 설정
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
