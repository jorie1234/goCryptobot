package main

import (
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"

	"github.com/urfave/cli/v2"
)

func InitCreateSellOrder() *cli.Command {

	return &cli.Command{
		Name:         "createsellorder",
		Aliases:      []string{"cso"},
		Usage:        "create a sell order. You can specify a corresponding buy order.",
		UsageText:    "If you specify a buy order the sell order will be associated with the buy order. If you dont specify a quantity, the executed quantity from the buy will be taken",
		Description:  "",
		ArgsUsage:    "",
		Category:     "",
		BashComplete: nil,
		Before:       nil,
		After:        nil,
		Action: func(c *cli.Context) error {
			binaClient := GetBinanceClient(c)
			quantity := c.String("quantity")
			price := c.String("price")
			buyorderid := c.Int64("buyorderid")
			symbol := c.String("symbol")
			takeQntFromBuy := c.Bool("takeqntfrombuy")

			var buyOrder *BinanceOrder

			if buyorderid > 0 {
				buyOrder = binaClient.GetOrderByID(buyorderid)
				if buyOrder == nil {
					fmt.Printf("BuyOrder %d could not be found !", buyorderid)
					return nil
				}
				if buyOrder.Order.Side != binance.SideTypeBuy {
					fmt.Printf("Order %d is not a buy order !", buyorderid)
					return nil
				}

				if buyOrder.Order.Status != binance.OrderStatusTypeFilled {
					fmt.Printf("Buy Order %d is not in status %s !", buyorderid, binance.OrderStatusTypeFilled)
					return nil
				}
				if len(symbol) == 0 {
					symbol = buyOrder.Order.Symbol
				}
			}
			if takeQntFromBuy {
				if buyOrder == nil {
					fmt.Printf("cannot find buy order woith id %d", buyorderid)
					return nil
				}
				quantity = buyOrder.Order.ExecutedQuantity
			}
			if buyorderid == 0 {
				if len(symbol) == 0 || len(price) == 0 {
					fmt.Printf("without BuyOrderID you have to use --symbol and --price and --quantity")
					return nil
				}
			}

			if len(quantity) == 0 {
				fmt.Printf("Error, no quantity specified. User --quantity or --takeqntfrombuy")
				return nil
			}
			pr, _ := strconv.ParseFloat(price, 8)
			quant, _ := strconv.ParseFloat(quantity, 8)
			fmt.Printf("create sell order for symbol %s price %s quant %s so sell will be at %.2f€ \n", symbol, price, quantity, pr*quant)

			if !Confirm("Perform sell ?") {
				return nil
			}

			sellOrder, err := binaClient.CreateSellOrder(symbol, price, quantity)
			if sellOrder == nil || err != nil {
				fmt.Printf("could not create sell order %s", err)
				return nil
			}

			binaClient.Store.SetSellForOrder(buyorderid, sellOrder.OrderID, 0.0)
			err = binaClient.Store.Save()
			if err != nil {
				fmt.Println(err)
				return nil
			}

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
				Required:    false,
				Hidden:      false,
				TakesFile:   false,
				Destination: nil,
				HasBeenSet:  false,
			},
			&cli.StringFlag{
				Name:        "buyorderid",
				Aliases:     []string{"bo"},
				Usage:       "id of buy order",
				EnvVars:     nil,
				FilePath:    "",
				Required:    false,
				Hidden:      false,
				TakesFile:   false,
				Destination: nil,
				HasBeenSet:  false,
			},
			&cli.StringFlag{
				Name:        "quantity",
				Aliases:     []string{"q"},
				Usage:       "quantity to sell",
				EnvVars:     nil,
				FilePath:    "",
				Required:    false,
				Hidden:      false,
				TakesFile:   false,
				Value:       "",
				Destination: nil,
				HasBeenSet:  false,
			},
			&cli.StringFlag{
				Name:        "price",
				Aliases:     []string{"p"},
				Usage:       "price to sell",
				EnvVars:     nil,
				FilePath:    "",
				Required:    false,
				Hidden:      false,
				TakesFile:   false,
				Value:       "",
				Destination: nil,
				HasBeenSet:  false,
			},
			&cli.BoolFlag{
				Name:        "takeqntfrombuy",
				Usage:       "Bool Flag to signal that the quantity should be taken from the buy order",
				EnvVars:     nil,
				FilePath:    "",
				Required:    false,
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
