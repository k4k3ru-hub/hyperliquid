//
// meta.go
//
package meta

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    myConstant "github.com/k4k3ru-hub/hyperliquid/go/constant"
    myRestDTO "github.com/k4k3ru-hub/hyperliquid/go/rest/dto"
)


const (
    TypeValue = "meta"
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

type Meta struct {
	Universe     []Universe         `json:"universe"`
	MarginTables []MarginTableEntry `json:"marginTables"`
}

type Universe struct {
	Name         string  `json:"name"`
	SzDecimals   int     `json:"szDecimals"`
	MaxLeverage  int     `json:"maxLeverage"`
	OnlyIsolated *bool   `json:"onlyIsolated,omitempty"`
	IsDelisted   *bool   `json:"isDelisted,omitempty"`
	MarginMode   *string `json:"marginMode,omitempty"`
}

type MarginTableEntry struct {
    ID    int         `json:"id"`
    Table MarginTable `json:"marginTable"`
}

type MarginTable struct {
	Description string       `json:"description"`
	MarginTiers []MarginTier `json:"marginTiers"`
}

type MarginTier struct {
	LowerBound  string `json:"lowerBound"`
	MaxLeverage int    `json:"maxLeverage"`
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
//   - 2026-04-05: Added.
//
func NewClient(parentClient ParentClient) (*Client, error) {
    // Guard.
    if parentClient == nil {
        return nil, fmt.Errorf("failed to create spot_meta client: missing required value: parent_client=null.")
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
// Send the request.
//
// Version:
//   - 2026-04-05: Added.
//
func (c *Client) Send(ctx context.Context) (*Meta, error) {
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
    var raw struct {
        Universe     []Universe          `json:"universe"`
        MarginTables [][]json.RawMessage `json:"marginTables"`
    }
    if err := json.Unmarshal(resBody, &raw); err != nil {
        return nil, err
    }

    // Build result.
    result := &Meta{
        Universe:     raw.Universe,
        MarginTables: make([]MarginTableEntry, 0, len(raw.MarginTables)),
    }

    // Parse marginTables entries.
    for i, entry := range raw.MarginTables {
        if len(entry) != 2 {
            return nil, fmt.Errorf("invalid marginTables[%d]: expected 2 elements, got %d", i, len(entry))
        }

        var id int
        if err := json.Unmarshal(entry[0], &id); err != nil {
            return nil, fmt.Errorf("failed to unmarshal marginTables[%d][0]: %w", i, err)
        }

        var table MarginTable
        if err := json.Unmarshal(entry[1], &table); err != nil {
            return nil, fmt.Errorf("failed to unmarshal marginTables[%d][1]: %w", i, err)
        }

        result.MarginTables = append(result.MarginTables, MarginTableEntry{
            ID:    id,
            Table: table,
        })
    }

    return result, nil
}
