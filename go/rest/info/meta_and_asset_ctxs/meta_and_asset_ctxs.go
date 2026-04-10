//
// meta_and_asset_ctxs.go
//
package meta_and_asset_ctxs

import (
    "context"
	"encoding/json"
    "fmt"
	"net/http"

    myConstant "github.com/k4k3ru-hub/hyperliquid/go/constant"
    myRestDTO "github.com/k4k3ru-hub/hyperliquid/go/rest/dto"
)


const (
	TypeValue = "metaAndAssetCtxs"
)


type MetaAndAssetCtxs struct {
	Universe []*UniverseEntry `json:"universe"`
	Assets   []*AssetEntry    `json:"assets"`
}
type Client struct {
    parent      ParentClient
    endpointURL string
    httpMethod  string
    reqBody     *myRestDTO.RequestBody
    httpHeader  http.Header
}
type UniverseEntry struct {
	Name         string `json:"name"`
	SzDecimals   int    `json:"szDecimals"`
	MaxLeverage  int    `json:"maxLeverage"`
	OnlyIsolated bool   `json:"onlyIsolated"`
}
type AssetEntry struct {
	DayNtlVlm    string   `json:"dayNtlVlm"`
	Funding      string   `json:"funding"`
	ImpactPxs    []string `json:"impactPxs"`
	MarkPx       string   `json:"markPx"`
	MidPx        string   `json:"midPx"`
	OpenInterest string   `json:"openInterest"`
	OraclePx     string   `json:"oraclePx"`
	Premium      string   `json:"premium"`
	PrevDayPx    string   `json:"prevDayPx"`
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


//
// Get the asset by the universe name.
//
func (object *MetaAndAssetCtxs) GetAssetByName(name string) *AssetEntry {
	findIndex := -1
	for i, entry := range object.Universe {
		if entry.Name == name {
			findIndex = i
			break
		}
	}
	if findIndex != -1 && findIndex <= len(object.Assets)-1 {
		return object.Assets[findIndex]
	}
	return nil
}


//
// Send a request.
//
func (c *Client) Send(ctx context.Context) (*MetaAndAssetCtxs, error) {
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
	result := &MetaAndAssetCtxs{}
	var rawData []json.RawMessage
	if err := json.Unmarshal(resBody, &rawData); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(rawData[0], result); err != nil {
		return nil, err
	}
	if len(rawData) > 1 {
		if err := json.Unmarshal(rawData[1], &result.Assets); err != nil {
			return nil, err
		}
	}

	return result, nil
}
