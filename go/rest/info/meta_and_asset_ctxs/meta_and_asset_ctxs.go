//
// meta_and_asset_ctxs.go
//
package meta_and_asset_ctxs

import (
	"encoding/json"
	"net/http"

	"github.com/k4k3ru-hub/hyperliquid/constant"
	"github.com/k4k3ru-hub/hyperliquid/rest"
)


const (
	TypeValue = "metaAndAssetCtxs"
)


type MetaAndAssetCtxs struct {
	Universe []*UniverseEntry `json:"universe"`
	Assets   []*AssetEntry    `json:"assets"`
}
type Client struct {
	*rest.Client
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


//
// New Client.
//
func NewClient() *Client {
	// New Client.
	client := rest.NewClient()

	// Set Endpoint URL.
	client.EndpointUrl = constant.BaseUrlRest + constant.ApiEndpointInfo

	// Set HTTP method.
	client.HttpMethod = http.MethodPost

	// Set RequestBody.
	client.RequestBody = &rest.RequestBody{
        Type: TypeValue,
    }

	return &Client{
		Client: client,
	}
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
func (c *Client) Send() (*MetaAndAssetCtxs, error) {
	// Send a request.
	resBody, err := c.Client.Send()
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
