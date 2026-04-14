//
// main.go
//
package main

import (
	"fmt"

	myCliRestAllMids          "github.com/k4k3ru-hub/hyperliquid/go/cmd/cli/rest/all_mids"
	myCliRestMeta             "github.com/k4k3ru-hub/hyperliquid/go/cmd/cli/rest/meta"
	myCliRestMetaAndAssetCtxs "github.com/k4k3ru-hub/hyperliquid/go/cmd/cli/rest/meta_and_asset_ctxs"
	myCliRestSpotMeta         "github.com/k4k3ru-hub/hyperliquid/go/cmd/cli/rest/spot_meta"
	spotMetaAndAssetCtxs      "github.com/k4k3ru-hub/hyperliquid/go/cmd/cli/rest/spot_meta_and_asset_ctxs"
	userFees                  "github.com/k4k3ru-hub/hyperliquid/go/cmd/cli/rest/user_fees"
    "github.com/k4k3ru-hub/hyperliquid/go/cmd/cli/websocket/l2book"

	"github.com/k4k3ru-hub/cli-go"
)


const (
	RestCommandName = "rest"
	RestCommandUsage = "REST API commands."
	WSCommandName = "ws"
    WSCommandUsage = "WS API commands."
)


//
// Main.
//
func main() {
	// Initialize CLI.
	myCli := cli.NewCli(nil)
	myCli.SetVersion("1.0.0")
	myCli.Command.SetDefaultConfigOption()

	// Add `rest` command.
	restCommand := cli.NewCommand(RestCommandName)
	restCommand.Usage = RestCommandUsage
	myCli.Command.Commands = append(myCli.Command.Commands, restCommand)

	// Add `ws` command.
	wsCommand := cli.NewCommand(WSCommandName)
	wsCommand.Usage = WSCommandUsage
	myCli.Command.Commands = append(myCli.Command.Commands, wsCommand)

	// Add `rest allMids` command.
	myCliRestAllMids.SetCommand(restCommand)

	// Add `rest meta` command.
	myCliRestMeta.SetCommand(restCommand)

	// Add `rest metaAndAssetCtxs` command.
	myCliRestMetaAndAssetCtxs.SetCommand(restCommand)

    // Add `rest spotMeta` command.
    myCliRestSpotMeta.SetCommand(restCommand)

    // Add `rest spotMetaAndAssetCtxs` command.
    spotMetaAndAssetCtxs.SetCommand(restCommand)

    // Add `rest userFees` command.
    userFees.SetCommand(restCommand)

    // Add `ws l2book` command.
    l2book.SetCommand(wsCommand)

	// Run the CLI.
    myCli.Run()
}


//
// Run.
//
func run(options map[string]*cli.Option) {
	// [TODO] Here is not supported yet.
	fmt.Printf("Started run function.\n")
}
