package ihs

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type simgosClient struct {
	domain     string
	httpClient *http.Client
}

func newSimgosClient(host string) *simgosClient {
	return &simgosClient{
		domain: "http://" + strings.TrimRight(host, "/"),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *simgosClient) sendIhs(ctx context.Context) (map[string]any, error) {
	url := c.domain + "/webservice/kemkes/ihs/patient/postIhs"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
