//
// spot_meta_and_asset_ctxs.go
//
package spot_meta_and_asset_ctxs

import (
    "context"
    "fmt"
    "strconv"
    "strings"

    myRest "github.com/k4k3ru-hub/hyperliquid/go/rest"

    "github.com/k4k3ru-hub/cli-go"
)


const (
    OptionNameToken = "token"
    OptionAliasToken = "t"

    ReqBodyType = "spotMetaAndAssetCtxs"
)


//
// Run.
//
func Run(options map[string]*cli.Option) {
    fmt.Printf("Started rest spotMetaAndAssetCtxs command.\n")

    // Create client.
    opt := myRest.DefaultClientOption()
    spotMetaAndAssetCtxsClient, err := myRest.NewClient(opt).InfoSpotMetaAndAssetCtxs()
    if err != nil {
        fmt.Printf("%s\n", err.Error())
        return
    }

    // Send API request.
    result, err := spotMetaAndAssetCtxsClient.Send(context.Background())
    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

    // Check the token option.
    var token string
    if tokenOption, ok := options[OptionNameToken]; ok {
        token = tokenOption.Value
    }

    tokenHeaders := []string{
	"Name",
	"SzDecimals",
	"WeiDecimals",
	"Index",
	"TokenID",
	"IsCanonical",
	"FullName",
}

var tokenData [][]interface{}
for _, tokenEntry := range result.Meta.Tokens {
	fullName := ""
	if tokenEntry.FullName != nil {
		fullName = *tokenEntry.FullName
	}

	rowData := []interface{}{
		tokenEntry.Name,
		strconv.Itoa(tokenEntry.SzDecimals),
		strconv.Itoa(tokenEntry.WeiDecimals),
		strconv.Itoa(tokenEntry.Index),
		tokenEntry.TokenID,
		strconv.FormatBool(tokenEntry.IsCanonical),
		fullName,
	}
	tokenData = append(tokenData, rowData)
}

cli.OutputTable(tokenHeaders, tokenData)



    // Set data.
    headers := []string{
        "Name",
        "Tokens",
        "Index",
        "IsCanonical",
        "DayNtlVlm",
        "MarkPx",
        "MidPx",
        "PrevDayPx",
    }

    var data [][]interface{}

    for _, universeEntry := range result.Meta.Universe {
        // Get symbol name.
        tokenName := universeEntry.Name
        tokenIndexes := make([]string, 0, len(universeEntry.Tokens))
        baseQuote := make([]string, 0, 2)
        for _, tokenIndex := range universeEntry.Tokens {
            tokenIndexes = append(tokenIndexes, strconv.Itoa(tokenIndex))

            if tokenIndex < len(result.Meta.Tokens) {
                baseQuote = append(baseQuote, result.Meta.Tokens[tokenIndex].Name)
            }
        }
        if len(baseQuote) == 2 {
            tokenName = baseQuote[0] + "/" + baseQuote[1]
        }

        // Filter after tokenName is resolved.
        if token != "" && strings.ToUpper(token) != strings.ToUpper(tokenName) {
            continue
        }

        // Guard.
        var dayNtlVlm string
        var markPx    string
        var midPx     string
        var prevDayPx string

        if universeEntry.Index < len(result.Prices) {
            p := result.Prices[universeEntry.Index]
            dayNtlVlm = p.DayNtlVlm
            markPx    = p.MarkPx
            midPx     = p.MidPx
            prevDayPx = p.PrevDayPx
        }

        rowData := []interface{}{
            tokenName,
            strings.Join(tokenIndexes, ","),
            strconv.Itoa(universeEntry.Index),
            strconv.FormatBool(universeEntry.IsCanonical),
            dayNtlVlm,
            markPx,
            midPx,
            prevDayPx,
        }

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
