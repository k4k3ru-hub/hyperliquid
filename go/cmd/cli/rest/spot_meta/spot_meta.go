//
// spot_meta.go
//
package spot_meta

import (
    "context"
    "fmt"
    "strings"

    myRest "github.com/k4k3ru-hub/hyperliquid/go/rest"

    "github.com/k4k3ru-hub/cli-go"
)


const (
    OptionNameToken = "token"
    OptionAliasToken = "t"

	ReqBodyType = "spotMeta"
)


//
// Run.
//
func Run(options map[string]*cli.Option) {
	fmt.Printf("Starting rest spotMeta command.\n")

    // Create client.
    opt := myRest.DefaultClientOption()
    restInfoSpotMetaClient, err := myRest.NewClient(opt).InfoSpotMeta()
    if err != nil {
        fmt.Printf("%s\n", err.Error())
        return
    }

	// Send API request.
	result, err := restInfoSpotMetaClient.Send(context.Background())
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

    // Check the token option.
    var token string
    if tokenOption, ok := options[OptionNameToken]; ok {
        token = tokenOption.Value
    }

	// Set data.
    fmt.Printf("[Tokens]\n")
	headers := []string{"name", "szDecimals", "weiDecimals", "index", "tokenId", "isCanonical", "evmContract", "fullName"}
	var data [][]interface{}
	for _, item := range result.Tokens {
		if token == "" || strings.ToUpper(token) == item.Name {
			data = append(data, []interface{}{
                item.Name,
                item.SzDecimals,
                item.WeiDecimals,
                item.Index,
                item.TokenID,
                item.IsCanonical,
                string(item.EVMContract),
                item.FullName,
            })
		}
	}

	// Output
	cli.OutputTable(headers, data)

    fmt.Printf("\n[Universe]\n")
    headers = []string{"name", "tokens", "index", "isCanonical"}
    data = [][]interface{}{}
    for _, item := range result.Universe {
        if token == "" || strings.ToUpper(token) == item.Name {
            data = append(data, []interface{}{
                item.Name,
                item.Tokens,
                item.Index,
                item.IsCanonical,
            })
        }
    }

    // Output
    cli.OutputTable(headers, data)
}


//
// Set command.
//
func SetCommand(parentCommand *cli.Command) {
	command := cli.NewCommand(ReqBodyType)
	parentCommand.Commands = append(parentCommand.Commands, command)
	command.Options[OptionNameToken] = &cli.Option{
		Alias: OptionAliasToken,
		HasValue: true,
	}
	command.Action = Run
}
