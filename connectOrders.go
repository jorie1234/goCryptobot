package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func InitConnectOrder() *cli.Command {

	return &cli.Command{
		Name:         "connectorder",
		Aliases:      []string{"co"},
		Usage:        "connect a buy order with a sell order ",
		UsageText:    "",
		Description:  "",
		ArgsUsage:    "",
		Category:     "",
		BashComplete: nil,
		Before:       nil,
		After:        nil,
		Action: func(c *cli.Context) error {
			binaClient := GetBinanceClient(c)

			buyID := c.Int64("buyorderid")
			sellID := c.Int64("sellorderid")

			err := binaClient.ConnectSellOrderWithBuyOrder(sellID, buyID)
			if err != nil {
				fmt.Printf("cannot connect order : %s", err)
				return nil
			}
			fmt.Printf("successfully connected buyorder %d with sellorder %d\n", buyID, sellID)
			_ = binaClient.Store.Save()
			return nil
		},
		OnUsageError: nil,
		Subcommands:  nil,
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:     "buyorderid",
				Aliases:  []string{"bid"},
				Usage:    "buy order id ",
				Required: true,
			},
			&cli.Int64Flag{
				Name:     "sellorderid",
				Aliases:  []string{"soi"},
				Usage:    "sell order id",
				Required: true,
			},
		},
	}
}
