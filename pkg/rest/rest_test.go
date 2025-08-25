package rest_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/devafterdark/project-jarvis/pkg/rest"
)

type mockRoundTripper struct {
	ResponseFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.ResponseFunc(req)
}

func newMockClient(statusCode int) http.Client {
	return http.Client{
		Transport: &mockRoundTripper{
			ResponseFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: statusCode,
					Body:       io.NopCloser(strings.NewReader(http.StatusText(statusCode))),
					Header:     make(http.Header),
				}, nil
			},
		},
	}
}

func TestClientRequestAPI(t *testing.T) {
	testCases := []struct {
		desc    string
		domain  string
		method  string
		client  http.Client
		wantErr bool
	}{
		{
			desc:    "status code 200",
			domain:  "https://example.com",
			method:  "GET",
			client:  newMockClient(200),
			wantErr: false,
		},
		{
			desc:    "status code 201",
			domain:  "https://example.com",
			method:  "POST",
			client:  newMockClient(201),
			wantErr: false,
		},
		{
			desc:    "status code 199",
			domain:  "https://example.com",
			method:  "GET",
			client:  newMockClient(199),
			wantErr: true,
		},
		{
			desc:    "status code 300",
			domain:  "https://example.com",
			method:  "GET",
			client:  newMockClient(300),
			wantErr: true,
		},
		{
			desc:    "invalid method",
			domain:  "https://example.com",
			method:  "INVALID",
			client:  newMockClient(400),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			client := rest.NewClient(tc.domain)
			_, err := client.RequestAPI(
				context.Background(), tc.method, "",
				rest.WithHTTPClient(tc.client),
			)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
