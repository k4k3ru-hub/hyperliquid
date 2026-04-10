//
// l2book.go
//
package l2book

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"

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

    // Create client.
    opt := websocket.DefaultClientOption()
    l2bookClient, err := websocket.NewClient(opt).SubscriptionL2Book(dto.Coin(coin))
    if err != nil {
        fmt.Printf("%s\n", err.Error())
        return
    }

    // Create context.
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    // Subscribe.
    if err := l2bookClient.Subscribe(ctx, func(book *dto.WsBook) {
        fmt.Printf("%+v\n", book) 
    }); err != nil {
        fmt.Printf("%s\n", err.Error())
        return
    }

    <- ctx.Done()


//	// Set data.
//    fmt.Printf("[Tokens]\n")
//	headers := []string{"name", "szDecimals", "weiDecimals", "index", "tokenId", "isCanonical", "evmContract", "fullName"}
//	var data [][]interface{}
//	for _, item := range result.Tokens {
//		if token == "" || strings.ToUpper(token) == item.Name {
//			data = append(data, []interface{}{
//                item.Name,
//                item.SzDecimals,
//                item.WeiDecimals,
//                item.Index,
//                item.TokenID,
//                item.IsCanonical,
//                string(item.EVMContract),
//                item.FullName,
//            })
//		}
//	}
//
//    // Output
//    cli.OutputTable(headers, data)
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



