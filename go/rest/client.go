//
// client.go
//
package rest

import (
	"bytes"
	"context"
	"encoding/json"
    "fmt"
	"io"
	"net"
	"net/http"
	"time"

    dto                      "github.com/k4k3ru-hub/hyperliquid/go/rest/dto"
    infoAllMids              "github.com/k4k3ru-hub/hyperliquid/go/rest/info/all_mids"
    infoMeta                 "github.com/k4k3ru-hub/hyperliquid/go/rest/info/meta"
    infoMetaAndAssetCtxs     "github.com/k4k3ru-hub/hyperliquid/go/rest/info/meta_and_asset_ctxs"
    infoSpotMeta             "github.com/k4k3ru-hub/hyperliquid/go/rest/info/spot_meta"
    infoSpotMetaAndAssetCtxs "github.com/k4k3ru-hub/hyperliquid/go/rest/info/spot_meta_and_asset_ctxs"
    userFees                 "github.com/k4k3ru-hub/hyperliquid/go/rest/info/user_fees"
)


type Client struct {
	httpClient *http.Client
	endpointURL string
	httpMethod string
    httpHeader http.Header
    body *dto.RequestBody
}

//
// Parameters:
//   - ConnectTimeout: Timeout for establishing the connection.
//
type ClientOption struct {
    ConnectTimeout time.Duration
}


//
// Get default client option.
//
// Version:
//   - 2026-04-04: Added.
//
func DefaultClientOption() *ClientOption {
    return &ClientOption{
        ConnectTimeout: 3 * time.Second,
    }
}

//
// New Client.
//
// Version:
//   - 2026-04-04: Added.
//
func NewClient(o *ClientOption) *Client {
    // Guard.
    if o == nil {
        o = DefaultClientOption()
    }

	return &Client{
        httpClient: &http.Client{
            Transport: &http.Transport{
                DialContext: (&net.Dialer{
                    Timeout: o.ConnectTimeout,
                }).DialContext,
            },
        },
	}
}


func (c *Client) InfoAllMids() (*infoAllMids.Client, error) {
    return infoAllMids.NewClient(c)
}


func (c *Client) InfoMeta() (*infoMeta.Client, error) {
    return infoMeta.NewClient(c)
}


func (c *Client) InfoMetaAndAssetCtxs() (*infoMetaAndAssetCtxs.Client, error) {
    return infoMetaAndAssetCtxs.NewClient(c)
}


func (c *Client) InfoSpotMeta() (*infoSpotMeta.Client, error) {
    return infoSpotMeta.NewClient(c)
}


func (c *Client) InfoSpotMetaAndAssetCtxs() (*infoSpotMetaAndAssetCtxs.Client, error) {
    return infoSpotMetaAndAssetCtxs.NewClient(c)
}


func (c *Client) InfoUserFees(user string) (*userFees.Client, error) {
    return userFees.NewClient(c, user)
}


func (c *Client) SetBody(body *dto.RequestBody) {
    c.body = body
}

func (c *Client) SetEndpointURL(endpointURL string) {
    c.endpointURL = endpointURL
}

func (c *Client) SetHttpMethod(method string) {
    c.httpMethod = method
}

func (c *Client) SetHttpHeader(header http.Header) {
    c.httpHeader = header
}


//
// Send a request.
//
// Version:
//   - 2026-04-04: Added.
//
func (c *Client) Send(ctx context.Context) ([]byte, error) {
    // Guard.
    if ctx == nil {
	    ctx = context.Background()
    }
    if c.endpointURL == "" {
        return nil, fmt.Errorf("failed to send request: missing required value: endpoint_url=null")
    }    
    if c.httpMethod == "" {
        return nil, fmt.Errorf("failed to send request: missing required value: http_method=null")
    }
    if len(c.httpHeader) == 0 {
        return nil, fmt.Errorf("failed to send request: missing required value: http_header=null")
    }

    endpointURL := c.endpointURL
    httpMethod := c.httpMethod

	// Set request body.
	var reqBody io.Reader
	if c.body != nil {
		byteBody, err := json.Marshal(*c.body)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}
		reqBody = bytes.NewBuffer(byteBody)
	}

	// Set Request.
	req, err := http.NewRequestWithContext(ctx, httpMethod, endpointURL, reqBody)
	if err != nil {
		return nil, err
	}

	// Set HTTP header.
    if len(c.httpHeader) > 0 {
        for k, v := range c.httpHeader {
            copied := make([]string, len(v))
            copy(copied, v)
            req.Header[k] = copied
        }
    }

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	return body, nil
}
