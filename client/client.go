package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	zotaBaseURL         = "https://api.zotapay-stage.com"
	absoluteUrlTemplate = "%s/%s"
)

// ZotaClientInterface allows easier testing
type ZotaClientInterface interface {
	Get(ctx context.Context, endpoint string, params url.Values) (*http.Response, error)
	Post(ctx context.Context, endpoint string, requestBody []byte) (*http.Response, error)
}

// ZotaClient is used to make API calls to ZOTA API (currently needed GET and POST)
type ZotaClient struct {
	client  *http.Client
	headers map[string]string
	baseUrl string
}

// GetZotaClient generates an authenticated ZOTA client
func GetZotaClient(secretKey string) ZotaClientInterface {
	return &ZotaClient{
		client: &http.Client{},
		headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": secretKey,
		},
		baseUrl: zotaBaseURL,
	}
}

func (c *ZotaClient) Get(ctx context.Context, endpoint string, params url.Values) (*http.Response, error) {
	qparams := ""
	if params != nil {
		qparams = params.Encode()
	}

	endpointUrl := fmt.Sprintf(absoluteUrlTemplate, c.baseUrl, endpoint)
	if qparams != "" {
		endpointUrl += "?" + qparams
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointUrl, nil)
	if err != nil {
		return nil, err
	}

	return c.performRequest(req)
}

func (c *ZotaClient) Post(ctx context.Context, endpoint string, requestBody []byte) (*http.Response, error) {
	endpointUrl := fmt.Sprintf(absoluteUrlTemplate, c.baseUrl, endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	return c.performRequest(req)
}

func (c *ZotaClient) performRequest(req *http.Request) (*http.Response, error) {
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return resp, err
		}
		return resp, fmt.Errorf("received non-2xx status code: %d, message: %v", resp.StatusCode, string(body))
	}

	return resp, nil
}
