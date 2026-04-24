//
// l2book.go
//
package l2book

import (
    "context"
    "encoding/json"
    "fmt"

    myWebsocketDTO "github.com/k4k3ru-hub/hyperliquid/go/websocket/dto"
)

type Client struct {
    parent ParentClient
    coin   myWebsocketDTO.Coin
}


type ParentClient interface {
    Subscribe(ctx context.Context, key string, payload []byte) error
    Unsubscribe(ctx context.Context, key string, payload []byte) error
}


func NewClient(parent ParentClient, coin myWebsocketDTO.Coin) (*Client, error) {
    // Guard.
    if parent == nil {
        return nil, fmt.Errorf("failed to create l2book client: missing required value: parent_client=null.")
    }

    return &Client{
        parent: parent,
        coin: coin,
    }, nil
}


func (c *Client) Subscribe(ctx context.Context) error {
    // Guard.
    if c == nil {
        return fmt.Errorf("failed to subscribe l2book: missing required value: client=null")
    }
    if c.coin == "" {
        return fmt.Errorf("failed to subscribe l2book: missing required value: coin=null")
    }
    if c.parent == nil {
        return fmt.Errorf("failed to subscribe l2book: missing required value: parent_client=null")
    }

    req := &myWebsocketDTO.SubscribeRequest{
        Method: myWebsocketDTO.MethodSubscribe,
        Subscription: myWebsocketDTO.Subscription{
            Type: myWebsocketDTO.SubscriptionTypeL2Book,
            Coin: string(c.coin),
        },
    }

    // Build subscription key.
    key, err := req.Subscription.BuildKey()
    if err != nil {
        return err
    }

    payload, err := json.Marshal(req)
    if err != nil {
        return fmt.Errorf("failed to send websocket json message: %w", err)
    }

    // Subscribe.
    return c.parent.Subscribe(ctx, key, payload)
}


func (c *Client) Unsubscribe(ctx context.Context) error {
    // Guard.
    if c == nil {
        return fmt.Errorf("failed to unsubscribe l2book: missing required value: client=null")
    }
    if c.coin == "" {
        return fmt.Errorf("failed to unsubscribe l2book: missing required value: coin=null")
    }
    if c.parent == nil {
        return fmt.Errorf("failed to unsubscribe l2book: missing required value: parent_client=null")
    }

    req := &myWebsocketDTO.SubscribeRequest{
        Method: myWebsocketDTO.MethodUnsubscribe,
        Subscription: myWebsocketDTO.Subscription{
            Type: myWebsocketDTO.SubscriptionTypeL2Book,
            Coin: string(c.coin),
        },
    }

    // Build subscription key.
    key, err := req.Subscription.BuildKey()
    if err != nil {
        return err
    }

    payload, err := json.Marshal(req)
    if err != nil {
        return fmt.Errorf("failed to send websocket json message: %w", err)
    }

    return c.parent.Unsubscribe(ctx, key, payload)
}
