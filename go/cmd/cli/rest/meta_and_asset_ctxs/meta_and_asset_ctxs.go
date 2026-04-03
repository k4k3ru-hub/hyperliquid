//
// meta_and_asset_ctxs.go
//
package meta_and_asset_ctxs

import (
	"fmt"
	"strconv"
	"strings"

	myRestInfoMetaAndAssetCtxs "github.com/k4k3ru-hub/hyperliquid/rest/info/meta_and_asset_ctxs"

	"github.com/k4k3ru-hub/cli-go"
)


const (
	OptionNameToken = "token"
	OptionAliasToken = "t"

	ReqBodyType = "metaAndAssetCtxs"
)


//
// Run.
//
func Run(options map[string]*cli.Option) {
	fmt.Printf("Started rest metaAndAssetCtxs command.\n")

	// Send API request.
	c := myRestInfoMetaAndAssetCtxs.NewClient()
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
	headers := []string{
		"Name",
		"SzDecimals",
		"MaxLeverage",
		"OnlyIsolated",
		"DayNtlVlm",
		"Funding",
		"ImpactPxs",
		"MarkPx",
		"MidPx",
		"OpenInterest",
		"OraclePx",
		"Premium",
		"PrevDayPx",
	}
	var data [][]interface{}
	for i, universeEntry := range result.Universe {
		if token != "" && strings.ToUpper(token) != universeEntry.Name && i+1 <= len(result.Assets) {
			continue
		}
		assetEntry := result.Assets[i]
		rowData := []interface{}{
			universeEntry.Name,
			strconv.Itoa(universeEntry.SzDecimals),
			strconv.Itoa(universeEntry.MaxLeverage),
			strconv.FormatBool(universeEntry.OnlyIsolated),
			assetEntry.DayNtlVlm,
			assetEntry.Funding,
			strings.Join(assetEntry.ImpactPxs, " "),
			assetEntry.MarkPx,
			assetEntry.MidPx,
			assetEntry.OpenInterest,
			assetEntry.OraclePx,
			assetEntry.Premium,
			assetEntry.PrevDayPx}

		data = append(data, rowData)
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

