package main

import (
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/urfave/cli/v2"
)

func InitReplaceSellOrder() *cli.Command {

	return &cli.Command{
		Name:         "replacesellorder",
		Aliases:      []string{"rso"},
		Usage:        "replace a sell order: delete sell order and create a new sell order",
		UsageText:    "",
		Description:  "",
		ArgsUsage:    "",
		Category:     "",
		BashComplete: nil,
		Before:       nil,
		After:        nil,
		Action: func(c *cli.Context) error {
			binaClient := GetBinanceClient(c)
			mult := c.Float64("mult")

			orderid := c.Int64("orderid")

			order := binaClient.GetOrderByID(orderid)
			if order == nil {
				fmt.Printf("Order %d could not be found !", orderid)
				return nil
			}
			if order.Order.Side != binance.SideTypeSell {
				fmt.Printf("Order %d is not a sell order !", orderid)
				return nil
			}

			if order.Order.Status != binance.OrderStatusTypeNew {
				fmt.Printf("Order %d is not in status %s !", orderid, binance.OrderStatusTypeNew)
				return nil
			}
			buyOrder := binaClient.GetBuyOrderIDforOrder(order)
			if buyOrder == nil {
				fmt.Printf("Order %d has no buy order ?!?!?!", orderid)
				return nil
			}
			fmt.Println("Sell Order:")
			fmt.Println(order.String())
			fmt.Println("Buy Order:")
			fmt.Println(buyOrder.String())
			if !Confirm("are these orders correct ?") {
				return nil
			}
			fmt.Println("cancel sell order...")
			err := binaClient.CancelOrder(order.Order.OrderID)
			if err != nil {
				fmt.Printf("Cannot cancel order %s", err)
				return nil
			}
			fmt.Printf("create new sell order with mult %f \n", mult)
			_, err = binaClient.PostSellForOrder(buyOrder.Order.OrderID, buyOrder.Order.Symbol, mult)
			if err != nil {
				fmt.Printf("error creating new sell order %sn", err)
				return nil
			}

			return nil
		},
		OnUsageError: nil,
		Subcommands:  nil,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "orderid",
				Aliases:     []string{"o"},
				Usage:       "id of sell order",
				EnvVars:     nil,
				FilePath:    "",
				Required:    true,
				Hidden:      false,
				TakesFile:   false,
				Destination: nil,
				HasBeenSet:  false,
			},
			&cli.StringFlag{
				Name:        "mult",
				Aliases:     []string{"m"},
				Usage:       "multiplier, replace sell order for CummulativeQuoteQuantity * mult",
				EnvVars:     nil,
				FilePath:    "",
				Required:    true,
				Hidden:      false,
				TakesFile:   false,
				Value:       "",
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
