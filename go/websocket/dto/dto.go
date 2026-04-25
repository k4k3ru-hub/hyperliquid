//
// dto.go
//
package dto

import (
    "encoding/json"
    "fmt"
    "strings"
)

const (
    MethodPing = "ping"
    MethodSubscribe = "subscribe"
    MethodUnsubscribe = "unsubscribe"

    SubscriptionTypeL2Book               = "l2Book"
    SubscriptionTypeNotification         = "notification"
    SubscriptionTypeWebData2             = "webData2"
    SubscriptionTypeSubscriptionResponse = "subscriptionResponse"
)

type Coin string
const (
    CoinBTC Coin = "BTC"
)

type SubscriptionKeyBuilder struct {
    Channel string `json:"channel"`
    Coin    string `json:"coin"`
    User    string `json:"user"`
    DEX     string `json:"dex"`
}

func (b *SubscriptionKeyBuilder) Build() string {
    // Guard.
    if b == nil || b.Channel == "" {
        return ""
    }

    parts := []string{b.Channel}

    if b.Coin != "" {
        parts = append(parts, b.Coin)
    }
    if b.User != "" {
        parts = append(parts, b.User)
    }
    if b.DEX != "" {
        parts = append(parts, b.DEX)
    }

    return strings.Join(parts, ":")
}

//
// Envelope is the common outer format for websocket responses.
//
// Examples:
//   - subscription ack:
//       { "channel": "subscriptionResponse", "data": { ... } }
//   - event push:
//       { "channel": "l2Book", "data": { ... } }
//   - post response:
//       { "channel": "post", "data": { "id": 123, "response": { ... } } }
//   - heartbeat response:
//       { "channel": "pong" }
//
// Version:
//   - 2026-04-06: Added.
//
type Envelope struct {
    Channel string          `json:"channel"`
    Data    json.RawMessage `json:"data,omitempty"`
}

func (e *Envelope) BuildKey() (string, error) {
    // Guard.
    if e == nil {
        return "", fmt.Errorf("failed to build key: missing requires value: receiver=null")
    }
    if e.Channel == "" {
        return "", fmt.Errorf("failed to build key: missing requires value: channel=null")
    }

    b := &SubscriptionKeyBuilder{
        Channel: e.Channel,
    }

    if len(e.Data) > 0 {
        if err := json.Unmarshal(e.Data, b); err != nil {
            return "", fmt.Errorf("[error] failed to build key: %w", err)
        }
    }

    return b.Build(), nil
}


type SubscribeRequest struct {
    Method       string       `json:"method"`
    Subscription Subscription `json:"subscription,omitempty"`
}

type Subscription struct {
    Type string `json:"type"`
    Coin string `json:"coin,omitempty"`
    User string `json:"user,omitempty"`
    DEX  string `json:"dex,omitempty"`
}

func (s *Subscription) BuildKey() (string, error) {
    // Guard.
    if s == nil {
        return "", fmt.Errorf("[error] failed to build key: missing requires value: receiver=null")
    }
    if s.Type == "" {
        return "", fmt.Errorf("[error] failed to build key: missing requires value: type=null")
    }

    b := &SubscriptionKeyBuilder{
        Channel: s.Type,
        Coin:    s.Coin,
        User:    s.User,
        DEX:     s.DEX,
    }

    return b.Build(), nil
}


type WsBookRaw struct {
    Coin   string      `json:"coin"`
    Levels [][]WsLevel `json:"levels"`
    Time   uint64      `json:"time"`
}

type WsBook struct {
    Coin string    `json:"coin"`
    Bids []WsLevel `json:"bids"`
	Asks []WsLevel `json:"asks"`
    Time uint64    `json:"time"`
}

type WsLevel struct {
    Price     string `json:"px"`
    Size      string `json:"sz"`
    NumOrders int    `json:"n"`
}

