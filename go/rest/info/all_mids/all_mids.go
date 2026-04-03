//
// all_mids.go
//
package all_mids

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/k4k3ru-hub/hyperliquid/constant"
	"github.com/k4k3ru-hub/hyperliquid/rest"
)


const (
	TypeValue = "allMids"
)

type AllMid struct {
	Token string `json:"token"`
	Mid   string `json:"mid"`
}
type Client struct {
    *rest.Client
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
// Send a request.
//
func (c *Client) Send() ([]*AllMid, error) {
    // Send a request.
    resBody, err := c.Client.Send()
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
