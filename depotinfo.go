package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v2"
)

func daySeconds(t time.Time) int {
	year, month, day := t.Date()
	t2 := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return int(t.Sub(t2).Seconds())
}

func InitDepotInfo() *cli.Command {

	return &cli.Command{
		Name:         "depotinfo",
		Aliases:      []string{"di"},
		Usage:        "list your new or filled and not selled orders",
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
			symb := strings.Split(symbol, ",")
			if last == "today" {
				last = fmt.Sprintf("%ds", daySeconds(time.Now()))
			}
			lastDuration, err := time.ParseDuration(last)
			if err != nil {
				fmt.Printf("Error parsing Last Parameter %s %v\n", last, err)
				return err
			}
			cumProfit := 0.0
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
					//"OrigQuantity",
					"ExQnt",
					"CumQuoteQnt",
					//"Value now",
					"Profit",
					"Side",
					"Status",
					"Time",
					"SellStat",
					//"Type",
					"SellOrder",
					"SellPrice",
					"SellTime",
				})
				t.SetStyle(table.StyleBold)
				i := 1
				unsoldQuantity := 0.0
				totalBuyQuantity := 0.0
				totalBuyNewQuantity := 0.0
				totalSellQuantity := 0.0
				totalSellNewQuantity := 0.0
				for _, o := range binaClient.Store.Orders {
					if o.Order.Symbol != v {
						continue
					}
					if o.Order.Status == binance.OrderStatusTypeCanceled {
						continue
					}

					exQuantity, _ := strconv.ParseFloat(o.Order.ExecutedQuantity, 8)
					avgPrc, _ := strconv.ParseFloat(price.Price, 8)
					cumQuote, _ := strconv.ParseFloat(o.Order.CummulativeQuoteQuantity, 8)
					profit := 0.0
					cumProfit = cumProfit + profit

					if o.Order.Status == binance.OrderStatusTypeFilled && o.Order.Side == binance.SideTypeBuy {
						totalBuyQuantity += exQuantity
					}
					if o.Order.Status == binance.OrderStatusTypeFilled && o.Order.Side == binance.SideTypeSell {
						totalSellQuantity += exQuantity
					}
					if o.Order.Status == binance.OrderStatusTypeNew && o.Order.Side == binance.SideTypeSell {
						totalSellNewQuantity += exQuantity
					}
					if o.Order.Status == binance.OrderStatusTypeNew && o.Order.Side == binance.SideTypeBuy {
						totalBuyNewQuantity += exQuantity
					}

					var sellStatus binance.OrderStatusType
					var (
						sellOrder BinanceOrder
					)
					sellOrderTime := ""
					profit = avgPrc*exQuantity - cumQuote

					if binaClient.GetRawSellOrderID(&o) == 1 {
						continue
					}
					if o.Order.Status == binance.OrderStatusTypeFilled && o.Order.Side == binance.SideTypeBuy && binaClient.DoesASellOrderExistForThisOrder(&o) {
						//continue
						sellOrderId := binaClient.GetSellOrderIDforOrder(&o)
						sellOrder = *binaClient.GetOrderByID(sellOrderId)
						sellOrderTime = time.Unix(sellOrder.Order.Time/1000, 0).Format(time.Stamp)
						if time.Unix(sellOrder.Order.UpdateTime/1000, 0).After(time.Unix(sellOrder.Order.Time/1000, 0)) {
							sellOrderTime = time.Unix(sellOrder.Order.UpdateTime/1000, 0).Format(time.Stamp)
						}
						sellStatus = sellOrder.Order.Status
						if sellStatus == binance.OrderStatusTypeFilled {
							sellCumQuote, _ := strconv.ParseFloat(sellOrder.Order.CummulativeQuoteQuantity, 8)
							profit = sellCumQuote - cumQuote
						}
					} else {
						unsoldQuantity += exQuantity
					}

					if !(o.OrderTimeYoungerThan(lastDuration) || binaClient.IsRelationYoungerThan(o, lastDuration)) {
						if o.Order.Status == binance.OrderStatusTypeExpired {
							continue
						}
						if !(sellStatus == binance.OrderStatusTypeNew) && len(sellStatus) > 0 {
							continue
						}
					}

					//fmt.Printf("%+#v", o)
					profitString := fmt.Sprintf("%2.2s", "-")
					profitColor := text.FgWhite
					if profit > 0 {
						profitString = fmt.Sprintf("%3.2f", profit)
						profitColor = text.FgGreen
					}
					if profit < 0 {
						profitString = fmt.Sprintf("%3.2f", profit)
						profitColor = text.FgRed
					}
					if o.Order.Side == "SELL" {
						profitString = fmt.Sprintf("%8.8s", "-")
						if binaClient.GetBuyOrderIDforOrder(&o) != nil {
							continue
						}
					}

					orderStateColor := text.FgWhite
					switch o.Order.Status {
					case binance.OrderStatusTypeNew:
						orderStateColor = text.FgYellow
					case binance.OrderStatusTypeFilled:
						if sellOrder.Order.OrderID > 0 {
							orderStateColor = text.FgCyan
							if sellOrder.Order.Status == binance.OrderStatusTypeFilled {
								orderStateColor = text.FgGreen
							}
						}
						if sellOrder.Order.OrderID == 0 {
							orderStateColor = text.FgRed
						}
					}
					//t.SetStyle(table.StyleColoredBlackOnGreenWhite)
					t.SetRowPainter(func(row table.Row) text.Colors {
						Color := text.FgWhite
						switch o.Order.Status {
						case binance.OrderStatusTypeNew:
							//Color=text.FgYellow
						case binance.OrderStatusTypeFilled:
							if sellOrder.Order.OrderID > 0 {
								Color = text.FgCyan
								if sellOrder.Order.Status == binance.OrderStatusTypeFilled {
									Color = text.FgGreen
								}
							}
							if sellOrder.Order.OrderID == 0 {
								Color = text.FgRed
							}
						}
						var colors []text.Color
						for range row {
							colors = append(colors, Color)
						}
						return colors
					})
					pr := o.Order.Price
					if o.Order.Type == binance.OrderTypeMarket {
						pr = "Market"
					}
					eq := o.Order.ExecutedQuantity
					if eq == "0.00000000" {
						eq = o.Order.OrigQuantity
					}
					t.AppendRow([]interface{}{
						//i,
						//o.Order.Symbol,
						o.Order.OrderID,
						//	o.Order.ClientOrderID,
						pr,
						//o.Order.OrigQuantity,
						eq,
						o.Order.CummulativeQuoteQuantity,
						//price.Price,
						//fmt.Sprintf("%.8f", quantity*avgPrc),
						profitColor.Sprint(profitString),
						//	o.Order.TimeInForce,
						o.Order.Side,
						orderStateColor.Sprint(o.Order.Status),
						//	o.Order.StopPrice,
						// o.Order.IcebergQuantity,
						time.Unix(o.Order.UpdateTime/1000, 0).Format(time.Stamp),
						sellStatus,
						//o.Order.Type,
						sellOrder.Order.OrderID, //binaClient.GetSellOrderIDforOrder(&o),
						sellOrder.Order.Price,
						//fmt.Sprintf("%+#v", sellOrder),
						sellOrderTime,
						//o.Order.UpdateTime,
						//o.Order.IsWorking,
						//o.Order.IsIsolated,
					}, table.RowConfig{
						AutoMerge: false,
					},
					)
					t.SetStyle(table.StyleBold)
					i++
				}
				//t.AppendFooter(table.Row{"", "", "", "", "", "", "", "", "", fmt.Sprintf("%11.8f", cumProfit), "", ""})
				if t.Length() > 0 {
					t.Render()
				}
				if unsoldQuantity > 0 {
					fmt.Printf("Unsold quantity: %.8f Total buy qnt %.8f total sell qnt %.8f diff %.8f total buy new %.8f total sell new %.8f diff %.8f\n",
						unsoldQuantity,
						totalBuyQuantity,
						totalSellQuantity,
						totalBuyQuantity-totalSellQuantity,
						totalBuyNewQuantity,
						totalSellNewQuantity,
						totalBuyNewQuantity-totalSellNewQuantity)
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
				Required:    false,
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

func GetBinanceClient(c *cli.Context) *BinanceClient {
	binaAPIKey := c.String("apikey")
	binaSecretKey := c.String("apisecret")
	binaClient := NewBinanceClient(binaAPIKey, binaSecretKey)
	return binaClient
}
