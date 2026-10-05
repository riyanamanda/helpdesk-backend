package bpjs

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"resty.dev/v3"
)

type Client struct {
	client         *resty.Client
	baseURL        string
	consumerID     string
	consumerSecret string
	userKey        string
}

func NewClient(baseUrl, consumerID, consumerSecret, userKey string) *Client {
	rc := resty.New()

	rc.SetHeaders(map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	})

	return &Client{
		client:         rc,
		baseURL:        baseUrl,
		consumerID:     consumerID,
		consumerSecret: consumerSecret,
		userKey:        userKey,
	}
}

func (c *Client) GetPeserta(ctx context.Context, nik string) ([]byte, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := generateSignature(c.consumerID, c.consumerSecret, timestamp)

	endpoint := fmt.Sprintf("%s/Peserta/nik/%s/tglSEP/%s", c.baseURL, nik, time.Now().Format("2006-01-02"))

	var result BPJSResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeaders(map[string]string{
			"X-cons-id":   c.consumerID,
			"X-timestamp": timestamp,
			"X-signature": signature,
			"user_key":    c.userKey,
		}).
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, err
	}
	if resp.IsStatusFailure() {
		return nil, errors.New("HTTP status failure dari BPJS")
	}

	if result.MetaData.Code != "200" {
		return []byte("{}"), nil
	}

	decryptKey := c.consumerID + c.consumerSecret + timestamp

	plainJson, err := DecryptResponse(decryptKey, result.Response)
	if err != nil {
		return nil, err
	}

	return []byte(plainJson), nil
}
