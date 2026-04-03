//
// all_mids.go
//
package all_mids

import (
	"fmt"
	"strings"

	myRestInfoAllMids "github.com/k4k3ru-hub/hyperliquid/rest/info/all_mids"

	"github.com/k4k3ru-hub/cli-go"
)


const (
    OptionNameToken = "token"
    OptionAliasToken = "t"

	ReqBodyType = "allMids"
)


//
// Run.
//
func Run(options map[string]*cli.Option) {
	fmt.Printf("Started rest allMids command.\n")

	// Send API request.
	c := myRestInfoAllMids.NewClient()
	result, err := c.Send()
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
	headers := []string{"Token", "Mid"}
	var data [][]interface{}
	for _, allMid := range result {
		if token == "" || strings.ToUpper(token) == allMid.Token {
			data = append(data, []interface{}{allMid.Token, allMid.Mid})
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
