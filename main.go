package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	"github.com/davecgh/go-spew/spew"

	"github.com/jedib0t/go-pretty/v6/table"
)

var (
	version    string
	date       string
	commitHash string
)

func main() {
	_ = godotenv.Load()

	//	binaAPIKey := os.Getenv("API_KEY")
	//	binaSecretKey := os.Getenv("API_SECRET")

	app := &cli.App{
		Usage:   "query data from binance",
		Version: fmt.Sprintf("%s\n\t Build %s\n\t Commit %s", version, date, commitHash),
		Commands: []*cli.Command{
			{
				Name:    "sellorder",
				Aliases: []string{"so"},
				Usage:   "sell an existing order",
				Action: func(c *cli.Context) error {
					binaAPIKey := c.String("apikey")
					binaSecretKey := c.String("apisecret")
					order := c.Int64("order")
					symbol := c.String("symbol")
					mult := c.Float64("mult")
					binaClient := NewBinanceClient(binaAPIKey, binaSecretKey)

					sellOrder, err := binaClient.PostSellForOrder(order, symbol, mult, 0.0)
					if err != nil {
						fmt.Println(err)
						return nil
					}
					fmt.Printf("%+#v", sellOrder)
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "order",
						Aliases:  []string{"o"},
						Usage:    "orderID of order to sell)",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "symbol",
						Aliases:  []string{"sy"},
						Usage:    "order symbol like ",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "mult",
						Aliases:  []string{"m"},
						Usage:    "multiplier, so price to sell = order.Price * mult",
						Required: true,
					},
				},
			},
			{
				Name:    "listprices",
				Aliases: []string{"lp"},
				Usage:   "list the prices",
				Action: func(c *cli.Context) error {
					binaAPIKey := c.String("apikey")
					binaSecretKey := c.String("apisecret")
					binaClient := NewBinanceClient(binaAPIKey, binaSecretKey)
					client := binaClient.client

					prices, err := client.NewListPricesService().Do(context.Background())
					if err != nil {
						fmt.Println(err)
						return nil
					}
					for _, p := range prices {
						fmt.Println(p)
					}
					return nil
				},
			},
			{
				Name:    "listorders",
				Aliases: []string{"lo"},
				Usage:   "list your orders",
				Action: func(c *cli.Context) error {
					binaAPIKey := c.String("apikey")
					binaSecretKey := c.String("apisecret")
					binaClient := NewBinanceClient(binaAPIKey, binaSecretKey)
					//client := binaClient.client

					symbol := c.String("symbol")
					symb := strings.Split(symbol, ",")

					t := table.NewWriter()
					t.SetOutputMirror(os.Stdout)
					t.AppendHeader(table.Row{"#",
						"Symbol",
						"OrderID",
						//	"ClientOrderID",
						"Price",
						"OrigQuantity",
						"ExQnt",
						"CumQuoteQnt",
						"AvgPrice now",
						"Value now",
						"Profit",
						"Status",
						//	"TimeInForce",
						"Type",
						"Side",
						//	"StopPrice",
						//	"IcebergQuantity",
						"Time",
						//	"UpdateTime",
						//	"IsWorking",
						//	"IsIsolated"
					})
					cumProfit := 0.0
					for _, v := range symb {
						binaClient.ListOrders(v)
						//orders, err := client.NewListOrdersService().Symbol(v).
						//	Do(context.Background())
						//if err != nil {
						//	fmt.Println(err)
						//	return nil
						//}

						t.AppendSeparator()
						t.SetStyle(table.StyleBold)
						i := 1
						for _, o := range binaClient.Store.Orders {
							if o.Order.Symbol != v {
								continue
							}
							if strings.EqualFold(string(o.Order.Status), "status") || strings.EqualFold(c.String("status"), "all") {
								//fmt.Printf("%+#v", o)
								avgPrice, _ := binaClient.GetAveragePrice(o.Order.Symbol)
								quantity, _ := strconv.ParseFloat(o.Order.ExecutedQuantity, 8)
								avgPrc, _ := strconv.ParseFloat(avgPrice.Price, 8)
								cumQuote, _ := strconv.ParseFloat(o.Order.CummulativeQuoteQuantity, 8)
								profit := quantity*avgPrc - cumQuote
								cumProfit = cumProfit + profit

								profitString := fmt.Sprintf("%11.8f", profit)
								if o.Order.Side == "SELL" {
									profitString = fmt.Sprintf("%8.8s", "-")
								}
								t.AppendRow([]interface{}{
									i,
									o.Order.Symbol,
									o.Order.OrderID,
									//	o.Order.ClientOrderID,
									o.Order.Price,
									o.Order.OrigQuantity,
									o.Order.ExecutedQuantity,
									o.Order.CummulativeQuoteQuantity,
									avgPrice.Price,
									fmt.Sprintf("%.8f", quantity*avgPrc),
									profitString,
									o.Order.Status,
									//	o.Order.TimeInForce,
									o.Order.Type,
									o.Order.Side,
									//	o.Order.StopPrice,
									// o.Order.IcebergQuantity,
									time.Unix(o.Order.Time/1000, 0),
									//o.Order.UpdateTime,
									//o.Order.IsWorking,
									//o.Order.IsIsolated,
								})
								i++
							}
						}
					}
					t.AppendFooter(table.Row{"", "", "", "", "", "", "", "", "", fmt.Sprintf("%11.8f", cumProfit), "", ""})
					t.Render()
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "status",
						Aliases:     []string{"st"},
						Usage:       "only list orders of this status (canceled,filled, new...)",
						EnvVars:     nil,
						FilePath:    "",
						Required:    false,
						Hidden:      false,
						TakesFile:   false,
						Value:       "all",
						Destination: nil,
						HasBeenSet:  false,
					},
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
				},
			},
			{
				Name:    "account",
				Aliases: []string{"a"},
				Usage:   "show your account",
				Action: func(c *cli.Context) error {
					binaAPIKey := c.String("apikey")
					binaSecretKey := c.String("apisecret")
					binaClient := NewBinanceClient(binaAPIKey, binaSecretKey)
					client := binaClient.client

					res, err := client.NewGetAccountService().Do(context.Background())
					if err != nil {
						fmt.Println(err)
						return nil
					}
					fmt.Printf("Can Deposit: %v\n", res.CanDeposit)
					fmt.Printf("Can Trade: %v\n", res.CanTrade)
					fmt.Printf("Can withdraw: %v\n", res.CanWithdraw)
					fmt.Printf("BuyerCommission: %d\n", res.BuyerCommission)
					fmt.Printf("TakerCommission: %d\n", res.TakerCommission)
					fmt.Printf("MakerCommission: %d\n", res.MakerCommission)
					fmt.Printf("SellerCommission: %d\n", res.SellerCommission)
					fmt.Printf("Balances: \n")
					t := table.NewWriter()
					t.SetOutputMirror(os.Stdout)
					t.SetStyle(table.StyleBold)
					t.AppendHeader(table.Row{
						text.AlignCenter.Apply("Symbol", 5),
						text.AlignCenter.Apply("Locked", 15),
						text.AlignCenter.Apply("Free", 15),
						text.AlignCenter.Apply("Price €", 15),
					})
					quoteSum := 0.0
					for _, v := range res.Balances {
						locked, _ := strconv.ParseFloat(v.Locked, 8)
						free, _ := strconv.ParseFloat(v.Free, 8)

						if !(locked == 0.0 && free == 0.0) {
							quote := "-"
							price, err := binaClient.GetAveragePrice(fmt.Sprintf("%s%s", v.Asset, "EUR"))
							if err == nil {
								priceQuote, _ := strconv.ParseFloat(price.Price, 8)
								quote = fmt.Sprintf("%.2f", (locked+free)*priceQuote)
								quoteSum += (locked + free) * priceQuote
							}
							t.AppendRow([]interface{}{
								text.AlignCenter.Apply(v.Asset, 5),
								text.AlignRight.Apply(v.Locked, 15),
								text.AlignRight.Apply(v.Free, 15),
								text.AlignRight.Apply(quote, 15),
							})
						}
					}
					t.AppendFooter(table.Row{"", "", "", text.AlignRight.Apply(fmt.Sprintf("%.2f", quoteSum), 15)})
					t.Render()
					return nil
				},
			},
			{
				Name:    "showsingleorder",
				Aliases: []string{"sso"},
				Usage:   "show a single order in raw API data",
				Action: func(c *cli.Context) error {
					binaAPIKey := c.String("apikey")
					binaSecretKey := c.String("apisecret")
					binaClient := NewBinanceClient(binaAPIKey, binaSecretKey)
					client := binaClient.client
					orderID := c.Int64("order")
					symbol := c.String("symbol")

					res, err := client.NewGetOrderService().OrderID(orderID).Symbol(symbol).Do(context.Background())
					if err != nil {
						fmt.Println(err)
						return nil
					}
					spew.Dump(res)
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "order",
						Aliases:     []string{"o"},
						Usage:       "order ID to fetch",
						EnvVars:     nil,
						FilePath:    "",
						Required:    false,
						Hidden:      false,
						TakesFile:   false,
						Value:       "all",
						Destination: nil,
						HasBeenSet:  false,
					},
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
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "apikey",
				Aliases: []string{"ak"},
				Usage:   "API Key",
				EnvVars: []string{"API_KEY"},
			},
			&cli.StringFlag{
				Name:    "apisecret",
				Aliases: []string{"as"},
				Usage:   "API Secret",
				EnvVars: []string{"API_SECRET"},
			},
		},
		EnableBashCompletion:   false,
		HideHelp:               false,
		HideHelpCommand:        false,
		HideVersion:            false,
		BashComplete:           nil,
		Before:                 nil,
		After:                  nil,
		Action:                 nil,
		CommandNotFound:        nil,
		OnUsageError:           nil,
		Compiled:               time.Time{},
		Authors:                nil,
		Copyright:              "",
		Reader:                 nil,
		Writer:                 nil,
		ErrWriter:              nil,
		ExitErrHandler:         nil,
		Metadata:               nil,
		ExtraInfo:              nil,
		CustomAppHelpTemplate:  "",
		UseShortOptionHandling: false,
	}

	app.Commands = append(app.Commands, InitDepotInfo())
	app.Commands = append(app.Commands, InitSellBot())
	app.Commands = append(app.Commands, InitWatchBot())
	app.Commands = append(app.Commands, InitConnectOrder())
	app.Commands = append(app.Commands, InitCreateMarketBuyOrder())
	app.Commands = append(app.Commands, InitCreateLimitBuyOrder())
	app.Commands = append(app.Commands, InitCancelOrder())
	app.Commands = append(app.Commands, InitCreateSellOrder())
	app.Commands = append(app.Commands, InitReplaceSellOrder())
	app.Commands = append(app.Commands, InitExchangeInfo())
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Confirm(text string) bool {
	fmt.Printf("%s [y/n]: ", text)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	response = strings.ToLower(strings.TrimSpace(response))

	if response == "y" || response == "yes" {
		return true
	}
	return false
}
