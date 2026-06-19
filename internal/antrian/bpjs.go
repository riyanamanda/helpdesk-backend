package antrian

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const antrolBasePath = "webservice/registrasionline/bpjs"

type antrolClient struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

func newAntrolClient(domain, username, password string) *antrolClient {
	return &antrolClient{
		baseURL:  domain + antrolBasePath,
		username: username,
		password: password,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *antrolClient) getToken(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/getToken", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("x-username", c.username)
	req.Header.Set("x-password", c.password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Metadata struct {
			Code int `json:"code"`
		} `json:"metadata"`
		Response struct {
			Token string `json:"token"`
		} `json:"response"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Metadata.Code != 200 {
		return "", fmt.Errorf("antrol: failed to get token (code %d)", result.Metadata.Code)
	}

	return result.Response.Token, nil
}

func (c *antrolClient) checkIn(ctx context.Context, kodeBooking int64) error {
	token, err := c.getToken(ctx)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(map[string]any{
		"kodebooking": fmt.Sprintf("%d", kodeBooking),
		"waktu":       time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/checkInAntrian", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("X-token", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("antrol: check-in failed (status %d)", resp.StatusCode)
	}

	return nil
}
