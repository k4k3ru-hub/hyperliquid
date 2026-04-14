//
// spot_meta_and_asset_ctxs.go
//
package spot_meta_and_asset_ctxs

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    myConstant "github.com/k4k3ru-hub/hyperliquid/go/constant"
    myRestDTO "github.com/k4k3ru-hub/hyperliquid/go/rest/dto"
)


const (
    TypeValue = "spotMetaAndAssetCtxs"
)


//
// Client.
//
type Client struct {
    parent      ParentClient
    endpointURL string
    httpMethod  string
    reqBody     *myRestDTO.RequestBody
    httpHeader  http.Header
}

type SpotMetaAndAssetCtxs struct {
    Meta   Meta
    Prices []Price
}

type Meta struct {
    Tokens   []Token    `json:"tokens"`
    Universe []Universe `json:"universe"`
}

type Token struct {
    Name        string          `json:"name"`
    SzDecimals  int             `json:"szDecimals"`
    WeiDecimals int             `json:"weiDecimals"`
    Index       int             `json:"index"`
    TokenID     string          `json:"tokenId"`
    IsCanonical bool            `json:"isCanonical"`
    EVMContract json.RawMessage `json:"evmContract"`
    FullName    *string         `json:"fullName"`
}

type Universe struct {
    Name        string `json:"name"`
    Tokens      []int  `json:"tokens"`
    Index       int    `json:"index"`
    IsCanonical bool   `json:"isCanonical"`
}

type Price struct {
    DayNtlVlm string `json:"dayNtlVlm"`
    MarkPx    string `json:"markPx"`
    MidPx     string `json:"midPx"`
    PrevDayPx string `json:"prevDayPx"`
}

type ParentClient interface {
    SetEndpointURL(endpointURL string)
    SetHttpMethod(method string)
    SetHttpHeader(header http.Header)
    SetBody(body *myRestDTO.RequestBody)
    Send(context.Context) ([]byte, error)
}


//
// New Client.
//
// Version:
//   - 2026-04-12: Added.
//
func NewClient(parentClient ParentClient) (*Client, error) {
    // Guard.
    if parentClient == nil {
        return nil, fmt.Errorf("failed to create spot meta and asset ctxs client: missing required value: parent_client=null")
    }

    // Create request body.
    reqBody := &myRestDTO.RequestBody{
        Type: TypeValue,
    }

    // Create http header.
    httpHeader := http.Header{
        "Content-Type": {myConstant.ContentTypeJson},
    }

    return &Client{
        parent: parentClient,
        endpointURL: myConstant.BaseUrlRest + myConstant.ApiEndpointInfo,
        httpMethod: http.MethodPost,
        reqBody: reqBody,
        httpHeader: httpHeader,
    }, nil
}


//
// Send a request.
//
func (c *Client) Send(ctx context.Context) (*SpotMetaAndAssetCtxs, error) {
    // Set date to parent client.
    c.parent.SetEndpointURL(c.endpointURL)
    c.parent.SetHttpMethod(c.httpMethod)
    c.parent.SetHttpHeader(c.httpHeader)
    c.parent.SetBody(c.reqBody)

    // Send a request.
    resBody, err := c.parent.Send(ctx)
    if err != nil {
        return nil, err
    }

    // Parse JSON data.
    result := &SpotMetaAndAssetCtxs{}
    var rawData []json.RawMessage
    if err := json.Unmarshal(resBody, &rawData); err != nil {
        return nil, err
    }
    if err := json.Unmarshal(rawData[0], &result.Meta); err != nil {
        return nil, err
    }
    if len(rawData) > 1 {
        if err := json.Unmarshal(rawData[1], &result.Prices); err != nil {
            return nil, err
        }
    }

    return result, nil
}


