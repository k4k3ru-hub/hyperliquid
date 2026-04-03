//
// client.go
//
package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/k4k3ru-hub/hyperliquid/constant"
)


type Client struct {
	EndpointUrl string
	HttpClient *http.Client
	HttpMethod string
	*RequestBody
	TimeoutSecond time.Duration
}
type RequestBody struct {
    Type string `json:"type"`
	User string `json:"user,omitempty"`
}


//
// New Client.
//
func NewClient() *Client {
	return &Client{
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}


//
// Send a request.
//
func (c *Client) Send() ([]byte, error) {
	// Set context.
	ctx := context.Background()
	var cancel context.CancelFunc
	if c.TimeoutSecond != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), c.TimeoutSecond*time.Second)
		defer cancel()
	}

	// Set request body.
	var reqBody io.Reader
	if c.RequestBody != nil {
		byteBody, err := json.Marshal(*c.RequestBody)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(byteBody)
	}

	// Set Request.
	req, err := http.NewRequestWithContext(ctx, c.HttpMethod, c.EndpointUrl, reqBody)
	if err != nil {
		return nil, err
	}

	// Set `Content-Type` header.
	req.Header.Set("Content-Type", constant.ContentTypeJson)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
