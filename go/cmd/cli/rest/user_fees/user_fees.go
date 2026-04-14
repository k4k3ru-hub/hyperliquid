//
// user_fees.go
//
package user_fees

import (
    "context"
    "fmt"
//    "strings"

    myRest "github.com/k4k3ru-hub/hyperliquid/go/rest"

    "github.com/k4k3ru-hub/cli-go"
)


const (
    OptionNameUser = "user"
    OptionAliasUser = "u"

    ReqBodyType = "userFees"
)


var (
    command = cli.NewCommand(ReqBodyType)
)


//
// Run.
//
func Run(options map[string]*cli.Option) {
    fmt.Printf("Starting rest userFees command.\n")

    // Check the user option.
    var user string
    if userOption, ok := options[OptionNameUser]; ok {
        user = userOption.Value
    }
    if user == "" {
        command.ShowUsage()
        return
    }

    // Create client.
    opt := myRest.DefaultClientOption()
    restInfoSpotMetaClient, err := myRest.NewClient(opt).InfoUserFees(user)
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

    // Output user fees.
    fmt.Printf("[UserFees]\n")
    headers := []string{"UserCrossRate", "UserAddRate", "UserSpotCrossRate", "UserSpotAddRate", "ActiveReferralDiscount", "Trial", "FeeTrialReward", "NextTrialAvailableTimestamp"}
    var data [][]interface{}
    data = append(data, []interface{}{
        result.UserCrossRate,
        result.UserAddRate,
        result.UserSpotCrossRate,
        result.UserSpotAddRate,
        result.ActiveReferralDiscount,
        result.Trial,
        result.FeeTrialReward,
        result.NextTrialAvailableTimestamp,
    })
    cli.OutputTable(headers, data)

    // Output dailyUserVlm.
    fmt.Printf("\n[DailyUserVlm]\n")
    headers = []string{"date", "userCross", "userAdd", "exchange"}
    data = [][]interface{}{}
    for _, item := range result.DailyUserVlm {
        data = append(data, []interface{}{
            item.Date,
            item.UserCross,
            item.UserAdd,
            item.Exchange,
        })
    }
    cli.OutputTable(headers, data)

    // Output .
    

}


//
// Set command.
//
func SetCommand(parentCommand *cli.Command) {
    parentCommand.Commands = append(parentCommand.Commands, command)
    command.Options[OptionNameUser] = &cli.Option{
        Alias: OptionAliasUser,
        HasValue: true,
    }
    command.Action = Run
}
