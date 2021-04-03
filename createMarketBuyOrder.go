package main

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"
)

func InitCreateMarketBuyOrder() *cli.Command {

	return &cli.Command{
		Name:         "createmarketbuyorder",
		Aliases:      []string{"cmbo"},
		Usage:        "create a market buy order for a symbol",
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
			//price := c.Float64("price")
			symbol := c.String("symbol")

			p, err := binaClient.GetAveragePrice(symbol)
			if err != nil {
				fmt.Printf("cannot get price for symbol %s -> %s\n", symbol, err)
				return nil
			}

			avgPrice, _ := strconv.ParseFloat(p.Price, 8)
			fmt.Printf("Buy %f of %s at market (this will cost you %.2f) ", quantity, symbol, quantity*avgPrice)
			if !Confirm("Please confirm ") {
				fmt.Println("buy canceled")
				return nil
			}
			order, err := binaClient.CreateMarketBuyOrder(symbol, quantity)
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
