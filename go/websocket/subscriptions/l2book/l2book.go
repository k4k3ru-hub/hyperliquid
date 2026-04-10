//
// l2book.go
//
package l2book

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    myWebsocketDTO "github.com/k4k3ru-hub/hyperliquid/go/websocket/dto"
)

type Client struct {
    parent ParentClient
    coin   myWebsocketDTO.Coin
}


type ParentClient interface {
    Subscribe(context.Context, *myWebsocketDTO.SubscribeRequest, func([]byte)) error
    Unsubscribe(context.Context, *myWebsocketDTO.SubscribeRequest) error
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


func (c *Client) Subscribe(ctx context.Context, handler func(*myWebsocketDTO.WsBook)) error {
    if c.coin == "" {
        return fmt.Errorf("failed to subscribe l2book: missing required value: coin=null")
    }

    req := &myWebsocketDTO.SubscribeRequest{
        Method: myWebsocketDTO.MethodSubscribe,
        Subscription: myWebsocketDTO.Subscription{
            Type: myWebsocketDTO.SubscriptionTypeL2Book,
            Coin: string(c.coin),
        },
    }

    return c.parent.Subscribe(ctx, req, func(data []byte) {
        raw := &myWebsocketDTO.WsBookRaw{}
        if err := json.Unmarshal(data, raw); err != nil {
            log.Printf("[error] failed to receive l2book event: %w", err)
            return
        }
        book := &myWebsocketDTO.WsBook{
            Coin: raw.Coin,
            Time: raw.Time,
        }
        if len(raw.Levels) > 0 {
            book.Bids = raw.Levels[0]
        }
        if len(raw.Levels) > 1 {
            book.Asks = raw.Levels[1]
        }
        handler(book)
    })
}


