//
// main.go
//
package main

import (
	"fmt"

	myCliRestAllMids          "github.com/k4k3ru-hub/hyperliquid/go/cli/rest/all_mids"
	myCliRestMetaAndAssetCtxs "github.com/k4k3ru-hub/hyperliquid/go/cli/rest/meta_and_asset_ctxs"

	"github.com/k4k3ru-hub/cli-go"
)


const (
	RestCommandName = "rest"
	RestCommandUsage = "REST API commands."
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

	// Add `rest allMids` command.
	myCliRestAllMids.SetCommand(restCommand)

	// Add `rest metaAndAssetCtxs` command.
	myCliRestMetaAndAssetCtxs.SetCommand(restCommand)

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
