//
// l2book.go
//
package l2book

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/k4k3ru-hub/hyperliquid/go/constant"
    "github.com/k4k3ru-hub/hyperliquid/go/websocket"
    "github.com/k4k3ru-hub/hyperliquid/go/websocket/dto"

    "github.com/k4k3ru-hub/cli-go"
)


const (
    OptionNameCoin = "coin"
    OptionAliasCoin = "c"

	ReqBodyType = "l2book"
)

var (
    command = cli.NewCommand(ReqBodyType)
)


type sessionHandler struct {}


//
// Run.
//
func Run(options map[string]*cli.Option) {
	fmt.Printf("Starting ws l2book command.\n")

    // Get coin option.
    coinOpt := options[OptionNameCoin]
    if coinOpt == nil {
        command.ShowUsage()
        return
    }
    coin := coinOpt.Value
    if coin == "" {
        command.ShowUsage()
        return
    }

    // Create context.
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    // Create client.
    endpointURL := constant.BaseUrlWebsocket + constant.ApiEndpointWebsocket
    sessHandler := &sessionHandler{}
    opt := websocket.DefaultClientOption()
    wsClient, err := websocket.NewClient(ctx, endpointURL, sessHandler, opt)
    if err != nil {
        fmt.Printf("%s\n", err.Error())
        return
    }
    l2bookClient, err := wsClient.SubscriptionL2Book(dto.Coin(coin))
    if err != nil {
        fmt.Printf("%s\n", err.Error())
        return
    }

    // Subscribe.
    if err := l2bookClient.Subscribe(ctx); err != nil {
        fmt.Printf("%s\n", err.Error())
        return
    }

    <- ctx.Done()
}


//
// Set command.
//
func SetCommand(parentCommand *cli.Command) {
	parentCommand.Commands = append(parentCommand.Commands, command)
	command.Options[OptionNameCoin] = &cli.Option{
		Alias: OptionAliasCoin,
		HasValue: true,
	}
	command.Action = Run
}


func (s *sessionHandler) HandleMessage(sess websocket.SessionContext, message []byte) {
//    fmt.Printf("%+v\n", string(message))

    var envelope dto.Envelope
    if err := json.Unmarshal(message, &envelope); err != nil {
        log.Printf("[error] failed to receive l2book event: %w", err)
        return
    }

    var bookRaw dto.WsBookRaw
    if err := json.Unmarshal(envelope.Data, &bookRaw); err != nil {
        log.Printf("[error] failed to receive l2book event: %w", err)
        return
    }

    book := &dto.WsBook{
        Coin: bookRaw.Coin,
        Time: bookRaw.Time,
    }
    if len(bookRaw.Levels) > 0 {
        book.Bids = bookRaw.Levels[0]
    }
    if len(bookRaw.Levels) > 1 {
        book.Asks = bookRaw.Levels[1]
    }
    fmt.Printf("%+v\n", book)
}


func (s *sessionHandler) HandleClose(sess websocket.SessionContext) {
    fmt.Printf("%+v\n", sess)
}
