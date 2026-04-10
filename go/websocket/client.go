//
// client.go
//
package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
    "log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	gorillaWebsocket "github.com/gorilla/websocket"

	myConstant "github.com/k4k3ru-hub/hyperliquid/go/constant"
	myWebsocketDTO "github.com/k4k3ru-hub/hyperliquid/go/websocket/dto"
	myWebsocketL2Book "github.com/k4k3ru-hub/hyperliquid/go/websocket/subscriptions/l2book"
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
type ClientOption struct {
	ConnectTimeout       time.Duration
	HandshakeTimeout     time.Duration
	ReconnectInterval    time.Duration
	MaxReconnectInterval time.Duration
}

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
	endpointURL string
	httpHeader  http.Header
	dialer      *gorillaWebsocket.Dialer

	conn   *gorillaWebsocket.Conn
	connMu sync.RWMutex

	sendCh chan []byte
	doneCh chan struct{}

	handlers   map[string]func([]byte)
	handlersMu sync.RWMutex

	subscriptions   map[string]myWebsocketDTO.SubscribeRequest
	subscriptionsMu sync.RWMutex

	connectMu sync.Mutex

	closeOnce sync.Once
	closed    atomic.Bool

	reconnecting        atomic.Bool
	reconnectInterval   time.Duration
	maxReconnectInterval time.Duration
}

//
// Get default client option.
//
// Version:
//   - 2026-04-06: Added.
//
func DefaultClientOption() *ClientOption {
	return &ClientOption{
		ConnectTimeout:       3 * time.Second,
		HandshakeTimeout:     5 * time.Second,
		ReconnectInterval:    1 * time.Second,
		MaxReconnectInterval: 30 * time.Second,
	}
}

//
// New client.
//
// Version:
//   - 2026-04-06: Added.
//
func NewClient(o *ClientOption) *Client {
	if o == nil {
		o = DefaultClientOption()
	}

	return &Client{
		endpointURL: myConstant.BaseUrlWebsocket + myConstant.ApiEndpointWebsocket,
		httpHeader:  make(http.Header),
		dialer: &gorillaWebsocket.Dialer{
			HandshakeTimeout: o.HandshakeTimeout,
			NetDialContext: (&net.Dialer{
				Timeout: o.ConnectTimeout,
			}).DialContext,
		},
		sendCh:               make(chan []byte, sendQueueSize),
		doneCh:               make(chan struct{}),
		handlers:             make(map[string]func([]byte)),
		subscriptions:        make(map[string]myWebsocketDTO.SubscribeRequest),
		reconnectInterval:    o.ReconnectInterval,
		maxReconnectInterval: o.MaxReconnectInterval,
	}
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

func (c *Client) SetEndpointURL(endpointURL string) {
	c.endpointURL = endpointURL
}

func (c *Client) SetHttpHeader(header http.Header) {
	if len(header) == 0 {
		c.httpHeader = make(http.Header)
		return
	}

	cloned := make(http.Header, len(header))
	for k, v := range header {
		copied := make([]string, len(v))
		copy(copied, v)
		cloned[k] = copied
	}
	c.httpHeader = cloned
}

//
// Register handler by route key.
//
// Supported examples:
//   - subscriptionResponse
//   - post
//   - pong
//   - l2Book:BTC
//   - l2Book:ETH
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) RegisterHandler(routeKey string, handler func([]byte)) {
	if c == nil || routeKey == "" || handler == nil {
		return
	}

	c.handlersMu.Lock()
	defer c.handlersMu.Unlock()

	c.handlers[routeKey] = handler
}

//
// Unregister handler by route key.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) UnregisterHandler(routeKey string) {
	if c == nil || routeKey == "" {
		return
	}

	c.handlersMu.Lock()
	defer c.handlersMu.Unlock()

	delete(c.handlers, routeKey)
}

//
// Connect websocket.
//
// Notes:
//   - Safe to call repeatedly.
//   - Starts read/write loops bound to the newly established connection.
//   - On reconnect, only the connection is swapped; handlers/subscriptions remain.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) Connect(ctx context.Context) error {
    // Guard.
	if c == nil {
		return fmt.Errorf("failed to connect websocket: missing required value: receiver=null")
	}
	if c.closed.Load() {
		return errors.New("failed to connect websocket: client is already closed")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if c.endpointURL == "" {
		return fmt.Errorf("failed to connect websocket: missing required value: endpoint_url=null")
	}

	c.connectMu.Lock()
	defer c.connectMu.Unlock()

	// Double-check after lock.
	c.connMu.RLock()
	if c.conn != nil {
		c.connMu.RUnlock()
		return nil
	}
	c.connMu.RUnlock()

	header := make(http.Header, len(c.httpHeader))
	for k, v := range c.httpHeader {
		copied := make([]string, len(v))
		copy(copied, v)
		header[k] = copied
	}

	conn, _, err := c.dialer.DialContext(ctx, c.endpointURL, header)
	if err != nil {
		return fmt.Errorf("failed to connect websocket: %w", err)
	}

	conn.SetReadLimit(maxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(initialReadWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	c.connMu.Lock()
	oldConn := c.conn
	c.conn = conn
	c.connMu.Unlock()

	if oldConn != nil {
		_ = oldConn.Close()
	}

	go c.pumpRead(conn)
	go c.pumpWrite(conn)

	return nil
}

//
// Close websocket client permanently.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) Close() error {
	if c == nil {
		return fmt.Errorf("failed to close websocket: missing required value: receiver=null")
	}

	var err error
	c.closeOnce.Do(func() {
		c.closed.Store(true)
		close(c.doneCh)

		c.connMu.Lock()
		if c.conn != nil {
			err = c.conn.Close()
			c.conn = nil
		}
		c.connMu.Unlock()
	})

	return err
}

//
// Subscribe.
//
// This method performs lazy connect on the first call.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) Subscribe(ctx context.Context, req *myWebsocketDTO.SubscribeRequest, handler func([]byte)) error {
	if c == nil {
		return fmt.Errorf("failed to subscribe websocket: missing required value: receiver=null")
	}
	if req == nil {
		return fmt.Errorf("failed to subscribe websocket: missing required value: subscribe_request=null")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to subscribe websocket: %w", err)
	}

    // Build subscription key.
    key, err := req.Subscription.BuildKey()
    if err != nil {
        return fmt.Errorf("failed to unsubscribe websocket: %w", err)
    }

	c.subscriptionsMu.Lock()
	c.subscriptions[key] = *req
	c.subscriptionsMu.Unlock()

    c.handlersMu.Lock()
    c.handlers[key] = handler
    c.handlersMu.Unlock()

	if err := c.Connect(ctx); err != nil {
		return err
	}

	return c.Write(ctx, reqBytes)
}

//
// Unsubscribe.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) Unsubscribe(ctx context.Context, req *myWebsocketDTO.SubscribeRequest) error {
	if c == nil {
		return fmt.Errorf("failed to unsubscribe websocket: missing required value: receiver=null")
	}
	if req == nil {
		return fmt.Errorf("failed to unsubscribe websocket: missing required value: request=null")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe websocket: %w", err)
	}

    // Build subscription key.
    key, err := req.Subscription.BuildKey()
    if err != nil {
        return fmt.Errorf("failed to unsubscribe websocket: %w", err)
    }

    c.subscriptionsMu.Lock()
    delete(c.subscriptions, key)
    c.subscriptionsMu.Unlock()

    c.handlersMu.Lock()
    delete(c.handlers, key)
    c.handlersMu.Unlock()

	return c.Write(ctx, reqBytes)
}


//
// Write JSON message.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) WriteJSON(ctx context.Context, v any) error {
	if c == nil {
		return fmt.Errorf("failed to write websocket json: missing required value: receiver=null")
	}
	if v == nil {
		return fmt.Errorf("failed to write websocket json: missing required value: payload=null")
	}

	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to write websocket json: %w", err)
	}

	return c.Write(ctx, b)
}

//
// Write raw websocket message.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) Write(ctx context.Context, data []byte) error {
	if c == nil {
		return fmt.Errorf("failed to write websocket message: missing required value: receiver=null")
	}
	if len(data) == 0 {
		return fmt.Errorf("failed to write websocket message: missing required value: data=null")
	}
	if c.closed.Load() {
		return errors.New("websocket client is closed")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to write websocket message: %w", ctx.Err())
	case <-c.doneCh:
		return errors.New("websocket client is closed")
	case c.sendCh <- append([]byte(nil), data...):
		return nil
	default:
		return errors.New("websocket send queue is full")
	}
}

//
// pumpRead reads Envelope via ReadJSON and dispatches by route key.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) pumpRead(conn *gorillaWebsocket.Conn) {
	for {
		select {
		case <-c.doneCh:
			return
		default:
		}

		env := &myWebsocketDTO.Envelope{}
		if err := conn.ReadJSON(env); err != nil {
			c.handleDisconnect(conn)
			return
		}

        // Build subscription key.
        key, err := env.BuildKey()
        if err != nil {
            log.Printf("[error] %s\n", err.Error())
            continue
        }
		if key == "" {
			continue
		}

		c.handlersMu.RLock()
		handler := c.handlers[key]
		c.handlersMu.RUnlock()
		if handler == nil {
			continue
		}

		handler(env.Data)
	}
}

//
// pumpWrite writes queued messages and sends Hyperliquid app-level ping.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) pumpWrite(conn *gorillaWebsocket.Conn) {
	ticker := time.NewTicker(appPingPeriod)
	defer ticker.Stop()

	pingReq := AppPingRequest{
		Method: "ping",
	}

	for {
		select {
		case <-c.doneCh:
			return

		case msg := <-c.sendCh:
			if !c.isCurrentConn(conn) {
				continue
			}

			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(gorillaWebsocket.TextMessage, msg); err != nil {
				c.handleDisconnect(conn)
				return
			}

		case <-ticker.C:
			if !c.isCurrentConn(conn) {
				return
			}

			pingBytes, err := json.Marshal(pingReq)
			if err != nil {
				continue
			}

			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(gorillaWebsocket.TextMessage, pingBytes); err != nil {
				c.handleDisconnect(conn)
				return
			}
		}
	}
}

//
// handleDisconnect swaps out the broken connection and starts reconnect loop.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) handleDisconnect(conn *gorillaWebsocket.Conn) {
	if c == nil || c.closed.Load() {
		return
	}

	c.connMu.Lock()
	if c.conn == conn {
		c.conn = nil
	}
	c.connMu.Unlock()

	_ = conn.Close()

	c.startReconnectLoop()
}

//
// startReconnectLoop reconnects with exponential backoff and re-subscribes.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) startReconnectLoop() {
	if c.closed.Load() {
		return
	}
	if !c.reconnecting.CompareAndSwap(false, true) {
		return
	}

	go func() {
		defer c.reconnecting.Store(false)

		backoff := c.reconnectInterval
		if backoff <= 0 {
			backoff = time.Second
		}

		for {
			if c.closed.Load() {
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			err := c.Connect(ctx)
			cancel()
			if err == nil {
				_ = c.resubscribeAll(context.Background())
				return
			}

			timer := time.NewTimer(backoff)
			select {
			case <-c.doneCh:
				timer.Stop()
				return
			case <-timer.C:
			}

			backoff *= 2
			if c.maxReconnectInterval > 0 && backoff > c.maxReconnectInterval {
				backoff = c.maxReconnectInterval
			}
		}
	}()
}

//
// resubscribeAll re-sends all stored subscription requests.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) resubscribeAll(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	c.subscriptionsMu.RLock()
	reqs := make([]myWebsocketDTO.SubscribeRequest, 0, len(c.subscriptions))
	for _, req := range c.subscriptions {
		reqs = append(reqs, req)
	}
	c.subscriptionsMu.RUnlock()

	for _, req := range reqs {
        reqBytes, err := json.Marshal(req)
        if err != nil {
            return fmt.Errorf("failed to resubscribe: %w", err)
        }

		if err := c.Write(ctx, reqBytes); err != nil {
			return err
		}
	}

	return nil
}

//
// isCurrentConn reports whether the given conn is still the active one.
//
// Version:
//   - 2026-04-06: Added.
//
func (c *Client) isCurrentConn(conn *gorillaWebsocket.Conn) bool {
	c.connMu.RLock()
	defer c.connMu.RUnlock()
	return c.conn == conn
}


