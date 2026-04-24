//
// client.go
//
package websocket

import (
	"context"
    "fmt"
	"time"

	myWebsocketDTO "github.com/k4k3ru-hub/hyperliquid/go/websocket/dto"
	myWebsocketL2Book "github.com/k4k3ru-hub/hyperliquid/go/websocket/subscriptions/l2book"

    k4k3ruWebsocket "github.com/k4k3ru-hub/websocket/go"
)

const (
	writeWait       = 10 * time.Second
	pongWait        = 70 * time.Second
	initialReadWait = pongWait
	maxMessageSize  = 1024 * 1024 // 1MB
	sendQueueSize   = 256

	// Hyperliquid recommends sending { "method": "ping" } if the subscribed
	// channel may be quiet for 60 seconds or more.
	appPingPeriod = 25 * time.Second
)


//
// AppPingRequest is the Hyperliquid websocket heartbeat request.
//
// Version:
//   - 2026-04-06: Added.
//
type AppPingRequest struct {
	Method string `json:"method"`
}

//
// ClientOption.
//
// Version:
//   - 2026-04-06: Added.
//
type ClientOption = k4k3ruWebsocket.ClientOption


//
// Client.
//
// Handler lifecycle:
//   - Handlers remain registered across reconnects.
//   - Subscriptions are stored and automatically re-sent after reconnect.
//
// Version:
//   - 2026-04-06: Added.
//
type Client struct {
    wsClient *k4k3ruWebsocket.Client
}

type SessionHandler = k4k3ruWebsocket.SessionHandler

type SessionContext = k4k3ruWebsocket.SessionContext


//
// Get default client option.
//
// Version:
//   - 2026-04-06: Added.
//
func DefaultClientOption() *ClientOption {
    return k4k3ruWebsocket.DefaultClientOption()
}

//
// New client.
//
// Version:
//   - 2026-04-06: Added.
//
func NewClient(ctx context.Context, endpointURL string, h SessionHandler, o *ClientOption) (*Client, error) {
    if h == nil {
        return nil, fmt.Errorf("failed to create client: missing required parameter: session_handler=null")
    }
	if o == nil {
		o = DefaultClientOption()
	}

    // Create new websocket client.
    wsClient, err := k4k3ruWebsocket.NewClient(ctx, endpointURL, h, o)
    if err != nil {
        return nil, err
    }

    return &Client{
        wsClient: wsClient,
    }, nil
}


//
// Create l2Book subscription client.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) SubscriptionL2Book(coin myWebsocketDTO.Coin) (*myWebsocketL2Book.Client, error) {
	return myWebsocketL2Book.NewClient(c, coin)
}


func (c *Client) Subscribe(ctx context.Context, key string, payload []byte) error {
    return c.wsClient.Subscribe(ctx, key, payload)
}

func (c *Client) Unsubscribe(ctx context.Context, key string, payload []byte) error {
    return c.wsClient. Unsubscribe(ctx, key, payload)
}




