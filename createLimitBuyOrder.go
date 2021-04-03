package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func InitCreateLimitBuyOrder() *cli.Command {

	return &cli.Command{
		Name:         "createlimitbuyorder",
		Aliases:      []string{"clbo"},
		Usage:        "create a limit buy order for a symbol",
		UsageText:    "",
		Description:  "",
		ArgsUsage:    "",
		Category:     "",
		BashComplete: nil,
		Before:       nil,
		After:        nil,
		Action: func(c *cli.Context) error {
			binaClient := GetBinanceClient(c)
			quantity := c.Float64("quantity")
			limit := c.Float64("limit")
			symbol := c.String("symbol")

			fmt.Printf("Buy %f of %s with limit of %f (this will cost you %.2f) ", quantity, symbol, limit, quantity*limit)
			if !Confirm("Please confirm ") {
				fmt.Println("buy canceled")
				return nil
			}
			order, err := binaClient.CreateLimitBuyOrder(symbol, quantity, limit)
			if err != nil || order == nil {
				fmt.Printf("Buy Order could not be created %s", err)
				return nil
			}

			fmt.Printf("Buy successful, order is %+#v", order)

			return nil
		},
		OnUsageError: nil,
		Subcommands:  nil,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "symbol",
				Aliases:     []string{"s"},
				Usage:       "Symbol",
				EnvVars:     nil,
				FilePath:    "",
				Required:    true,
				Hidden:      false,
				TakesFile:   false,
				Destination: nil,
				HasBeenSet:  false,
			},
			&cli.Float64Flag{
				Name:        "quantity",
				Aliases:     []string{"q"},
				Usage:       "quantity to buy",
				EnvVars:     nil,
				FilePath:    "",
				Required:    true,
				Hidden:      false,
				Destination: nil,
				HasBeenSet:  false,
			},
			&cli.Float64Flag{
				Name:        "limit",
				Aliases:     []string{"l"},
				Usage:       "buy limit",
				EnvVars:     nil,
				FilePath:    "",
				Required:    true,
				Hidden:      false,
				Destination: nil,
				HasBeenSet:  false,
			},
		},
		SkipFlagParsing:        false,
		HideHelp:               false,
		HideHelpCommand:        false,
		Hidden:                 false,
		UseShortOptionHandling: false,
		HelpName:               "",
		CustomHelpTemplate:     "",
	}
}
