//
// all_mids.go
//
package all_mids

import (
    "context"
	"encoding/json"
    "fmt"
	"net/http"
	"sort"

	myConstant "github.com/k4k3ru-hub/hyperliquid/go/constant"
	myRestDTO "github.com/k4k3ru-hub/hyperliquid/go/rest/dto"
)


const (
	TypeValue = "allMids"
)

type Client struct {
    parent      ParentClient
    endpointURL string
    httpMethod  string
    httpHeader  http.Header
    reqBody     *myRestDTO.RequestBody
}

type AllMid struct {
	Token string `json:"token"`
	Mid   string `json:"mid"`
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
func NewClient(parentClient ParentClient) (*Client, error) {
    // Guard.
    if parentClient == nil {
        return nil, fmt.Errorf("failed to create client: missing required value: parent_client=null.")
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


func (c *Client) SetDEX(dex string) {
    c.reqBody.DEX = dex
}


//
// Send a request.
//
func (c *Client) Send(ctx context.Context) ([]*AllMid, error) {
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
    var result []*AllMid
	midsMap := make(map[string]string)
    if err := json.Unmarshal(resBody, &midsMap); err != nil {
        return nil, err
    }
	if len(midsMap) == 0 {
		return result, nil
	}
	var mapKeys []string
	for token, _ := range midsMap {
		mapKeys = append(mapKeys, token)
	}
	sort.Strings(mapKeys)
	for _, mapKey := range mapKeys {
		result = append(result, &AllMid{
			Token: mapKey,
			Mid: midsMap[mapKey],
		})
	}

    return result, nil
}
