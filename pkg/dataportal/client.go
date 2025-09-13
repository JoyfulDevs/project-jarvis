package dataportal

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/joyfuldevs/project-jarvis/pkg/rest"
)

const domain = "https://apis.data.go.kr"

type clientOptions struct {
	httpClient http.Client
}

type Option func(*clientOptions)

func WithHTTPClient(client http.Client) Option {
	return func(o *clientOptions) {
		o.httpClient = client
	}
}

type Client struct {
	AuthKey string
	options *clientOptions
}

func NewClient(authKey string, opts ...Option) *Client {
	options := &clientOptions{
		httpClient: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
					MaxVersion: tls.VersionTLS13,
					CipherSuites: []uint16{
						// for *.data.go.kr
						tls.TLS_RSA_WITH_AES_128_CBC_SHA,
						tls.TLS_RSA_WITH_AES_256_CBC_SHA,
					},
				},
			},
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(options)
	}
	return &Client{AuthKey: authKey, options: options}
}

func (c *Client) ListHolidays(
	ctx context.Context,
	req *ListHolidaysRequest,
) (*ListHolidaysResponse, error) {
	params := map[string]string{
		"serviceKey": c.AuthKey,
		"solYear":    strconv.Itoa(req.Year),
		"solMonth":   fmt.Sprintf("%02d", req.Month),
	}

	client := rest.NewClient(domain)
	data, err := client.RequestAPI(
		ctx,
		"GET",
		"/B090041/openapi/service/SpcdeInfoService/getRestDeInfo",
		rest.WithHTTPClient(c.options.httpClient),
		rest.WithParams(params),
	)
	if err != nil {
		return nil, err
	}

	resp := &ListHolidaysResponse{}
	if err := xml.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetUltraShortTermForecast(
	ctx context.Context,
	req *UltraShortTermForecastRequest,
) (*UltraShortTermForecastResponse, error) {
	params := map[string]string{
		"serviceKey": c.AuthKey,
		"dataType":   "JSON",
		"pageNo":     strconv.Itoa(req.Page),
		"numOfRows":  strconv.Itoa(req.Count),
		"base_date":  req.BaseDate,
		"base_time":  req.BaseTime,
		"nx":         strconv.Itoa(req.NX),
		"ny":         strconv.Itoa(req.NY),
	}

	client := rest.NewClient(domain)
	data, err := client.RequestAPI(
		ctx,
		"GET",
		"/1360000/VilageFcstInfoService_2.0/getUltraSrtFcst",
		rest.WithHTTPClient(c.options.httpClient),
		rest.WithParams(params),
	)
	if err != nil {
		return nil, err
	}

	resp := &UltraShortTermForecastResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
