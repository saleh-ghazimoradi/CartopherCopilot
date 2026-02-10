package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type RestClient struct {
	client       *http.Client
	baseURL      string
	defaultToken string
	useToken     bool
	logger       *slog.Logger
}

func (c *RestClient) prepareRequest(method, path string, queryParams map[string]string, body any) (*http.Request, error) {
	fullURL, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, err
	}

	if len(queryParams) > 0 {
		q := fullURL.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		fullURL.RawQuery = q.Encode()
	}

	var bodyReader *bytes.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	} else {
		bodyReader = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, fullURL.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.useToken && c.defaultToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.defaultToken)
	}

	return req, nil
}

func (c *RestClient) Get(path string, params map[string]string) ([]byte, error) {
	req, err := c.prepareRequest("GET", path, params, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	c.logger.Debug("REST API call", "method", req.Method, "url", req.URL.String(), "status", resp.StatusCode)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error from Get(): %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (c *RestClient) Post(path string, body any) ([]byte, error) {
	req, err := c.prepareRequest("POST", path, nil, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	c.logger.Debug("REST API call", "method", req.Method, "url", req.URL.String(), "status", resp.StatusCode)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error Post(): %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (c *RestClient) WithToken() *RestClient {
	clone := *c
	clone.useToken = true
	return &clone
}

func NewRestClient(baseURL, defaultToken string, logger *slog.Logger) *RestClient {
	client := &http.Client{}
	return &RestClient{
		client:       client,
		baseURL:      baseURL,
		defaultToken: defaultToken,
		logger:       logger,
	}
}
