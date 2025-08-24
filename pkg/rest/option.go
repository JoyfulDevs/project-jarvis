package rest

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

type options struct {
	httpClient http.Client
	params     string
	headers    map[string]string
	body       io.Reader
}

type Option func(*options)

func WithParams(params map[string]string) Option {
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}
	return func(o *options) {
		o.params = "?" + values.Encode()
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(o *options) {
		o.headers = headers
	}
}

func WithBody(body []byte) Option {
	return func(o *options) {
		o.body = bytes.NewReader(body)
	}
}

func WithHTTPClient(client http.Client) Option {
	return func(o *options) {
		o.httpClient = client
	}
}
