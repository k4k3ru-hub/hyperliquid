//
// meta.go
//
package meta

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

	ReqBodyType = "meta"
)


//
// Run.
//
func Run(options map[string]*cli.Option) {
	fmt.Printf("Starting rest meta command.\n")

    // Create client.
    opt := myRest.DefaultClientOption()
    restInfoMetaClient, err := myRest.NewClient(opt).InfoMeta()
    if err != nil {
        fmt.Printf("%s\n", err.Error())
        return
    }

	// Send API request.
	result, err := restInfoMetaClient.Send(context.Background())
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
    fmt.Printf("[Universe]\n")
	headers := []string{"name", "szDecimals", "maxLeverage"}
	var data [][]interface{}
	for _, item := range result.Universe {
		if token == "" || strings.ToUpper(token) == item.Name {
			data = append(data, []interface{}{
                item.Name,
                item.SzDecimals,
                item.MaxLeverage,
            })
		}
	}

	// Output
	cli.OutputTable(headers, data)

    fmt.Printf("\n[marginTables]\n")
    headers = []string{"id", "description", "marginTiers"}
    data = [][]interface{}{}
    for _, item := range result.MarginTables {
        data = append(data, []interface{}{
            item.ID,
            item.Table.Description,
            item.Table.MarginTiers,
        })
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
