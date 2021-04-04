package main

import (
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/urfave/cli/v2"
)

func InitCancelOrder() *cli.Command {

	return &cli.Command{
		Name:         "cancelorder",
		Aliases:      []string{"co"},
		Usage:        "cancle an order",
		UsageText:    "",
		Description:  "",
		ArgsUsage:    "",
		Category:     "",
		BashComplete: nil,
		Before:       nil,
		After:        nil,
		Action: func(c *cli.Context) error {
			binaClient := GetBinanceClient(c)

			orderid := c.Int64("orderid")

			order := binaClient.GetOrderByID(orderid)
			if order == nil {
				fmt.Printf("Order %d could not be found !", orderid)
				return nil
			}

			if order.Order.Status != binance.OrderStatusTypeNew {
				fmt.Printf("Order %d is not in status %s !", orderid, binance.OrderStatusTypeNew)
				return nil
			}
			fmt.Println("Order:")
			fmt.Println(order.String())
			if !Confirm("is this order correct ?") {
				return nil
			}
			err := binaClient.CancelOrder(order.Order.OrderID)
			if err != nil {
				fmt.Printf("Cannot cancel order %s", err)
				return nil
			}
			fmt.Printf("order %d canceld \n", orderid)

			return nil
		},
		OnUsageError: nil,
		Subcommands:  nil,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "orderid",
				Aliases:     []string{"o"},
				Usage:       "id of order to cancel",
				EnvVars:     nil,
				FilePath:    "",
				Required:    true,
				Hidden:      false,
				TakesFile:   false,
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
