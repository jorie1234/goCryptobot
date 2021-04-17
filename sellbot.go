package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
)

func InitSellBot() *cli.Command {

	return &cli.Command{
		Name:         "sellbot",
		Aliases:      []string{"sb"},
		Usage:        "let the sellbot sell all your open orders",
		UsageText:    "",
		Description:  "",
		ArgsUsage:    "",
		Category:     "",
		BashComplete: nil,
		Before:       nil,
		After:        nil,
		Action: func(c *cli.Context) error {
			binaClient := GetBinanceClient(c)

			symbol := c.String("symbol")
			last := c.String("last")
			mult := c.Float64("mult")
			symb := strings.Split(symbol, ",")
			lastDuration, err := time.ParseDuration(last)
			var ordersToSell []BinanceOrder

			if err != nil {
				fmt.Printf("Error parsing Last Parameter %s %v\n", last, err)
				return err
			}
			for _, v := range symb {
				binaClient.ListOrders(v)

				if len(v) == 0 {
					continue
				}

				price, err := binaClient.GetAveragePrice(v)
				if err != nil {
					fmt.Printf("Error getting Price %v\n", err)
					return err
				}

				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)

				fmt.Printf("Price %s\n", price.Price)

				t.AppendSeparator()
				t.AppendHeader(table.Row{
					"OrderID",
					"Price",
					"ExQnt",
					"CumQuoteQnt",
					"Side",
					"Status",
					"SellQuoteQnt",
					"SellPrice",
					"Time",
				})
				t.SetStyle(table.StyleBold)
				i := 1
				for _, o := range binaClient.Store.Orders {
					if o.Order.Symbol != v {
						continue
					}
					if o.Order.Side != binance.SideTypeBuy {
						continue
					}
					if o.Order.Status != binance.OrderStatusTypeFilled {
						continue
					}
					if time.Unix(o.Order.Time/1000, 0).Before(time.Now().Add(-lastDuration)) {
						continue
					}
					if binaClient.DoesASellOrderExistForThisOrder(&o) {
						continue
					}
					cumQuote, _ := strconv.ParseFloat(o.Order.CummulativeQuoteQuantity, 8)
					orderPrice, _ := strconv.ParseFloat(o.Order.Price, 8)
					ordersToSell = append(ordersToSell, o)
					t.AppendRow([]interface{}{
						o.Order.OrderID,
						o.Order.Price,
						o.Order.ExecutedQuantity,
						o.Order.CummulativeQuoteQuantity,
						o.Order.Side,
						o.Order.Status,
						fmt.Sprintf("%.8f", cumQuote*mult),
						fmt.Sprintf("%.8f", orderPrice*mult),
						time.Unix(o.Order.Time/1000, 0).Format(time.Stamp),
					}, table.RowConfig{
						AutoMerge: false,
					},
					)
					i++
				}
				//t.AppendFooter(table.Row{"", "", "", "", "", "", "", "", "", fmt.Sprintf("%11.8f", cumProfit), "", ""})
				t.Render()
				fmt.Printf("Refresh Exchangeinfo....")
				err = binaClient.GetExchangeInfo()
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return nil
				}
				fmt.Printf("done\n")
				for _, o := range ordersToSell {
					_, err := binaClient.PostSellForOrder(o.Order.OrderID, o.Order.Symbol, mult)
					if err != nil {
						fmt.Println(err)
					}
				}
			}

			return nil
		},
		OnUsageError: nil,
		Subcommands:  nil,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "symbol",
				Aliases:     []string{"sy"},
				Usage:       "order symbol like ",
				EnvVars:     nil,
				FilePath:    "",
				Required:    true,
				Hidden:      false,
				TakesFile:   false,
				Value:       "BTCEUR",
				Destination: nil,
				HasBeenSet:  false,
			},
			&cli.StringFlag{
				Name:        "last",
				Aliases:     []string{"l"},
				Usage:       "time period e.g. 24h for the last 24 hours",
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
				Name:        "mult",
				Aliases:     []string{"m"},
				Usage:       "multiplier, create sell order for CummulativeQuoteQuantity * mult",
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
