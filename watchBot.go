package main

import (
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
	"time"
)

func InitWatchBot() *cli.Command {

	return &cli.Command{
		Name:         "watchbot",
		Aliases:      []string{"wb"},
		Usage:        "the watchbot waits for new filled buy orders and send them to telegram",
		UsageText:    "",
		Description:  "",
		ArgsUsage:    "",
		Category:     "",
		BashComplete: nil,
		Before:       nil,
		After:        nil,
		Action: func(c *cli.Context) error {
			binaClient := GetBinanceClient(c)

			telegrambotkey := c.String("telegrambotkey")
			telegramchatid := c.String("telegramchatid")
			symbol := c.String("symbol")
			repeat := c.String("repeat")
			last := c.String("last")
			symb := strings.Split(symbol, ",")
			lastDuration, err := time.ParseDuration(last)
			if err != nil {
				fmt.Printf("Error parsing Last Parameter %s %v\n", last, err)
				return err
			}
			var repeatDuration time.Duration
			if len(repeat)>0 {
				repeatDuration, err = time.ParseDuration(repeat)
				if err != nil {
					fmt.Printf("Error parsing repeat Parameter %s %v\n", last, err)
					return err
				}
			}
			for {
				binaClient.Store.Load()
				dataChanged:=false
				for _, v := range symb {
					binaClient.ListOrders(v)

					if len(v) == 0 {
						continue
					}

					t := table.NewWriter()
					t.SetOutputMirror(os.Stdout)

					t.AppendSeparator()
					t.AppendHeader(table.Row{
						"OrderID",
						"Price",
						"ExQnt",
						"CumQuoteQnt",
						"Side",
						"Status",
						"Time",
					})
					t.SetStyle(table.StyleBold)
					i := 1
					for k, o := range binaClient.Store.Orders {
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
						if o.FilledStatusSendToTelegram == true {
							continue
						}
						if len(telegrambotkey) > 0 {
							telegram := NewTelegram(telegrambotkey, telegramchatid)
							telegram.SendTelegramMessage(fmt.Sprintf("Bought %s at %s ExQnt %s, CumQuoteQnt %s€",
								o.Order.Symbol,
								o.Order.Price,
								o.Order.ExecutedQuantity,
								o.Order.CummulativeQuoteQuantity,
							))
							binaClient.Store.Orders[k].FilledStatusSendToTelegram = true
							dataChanged=true
						}
						t.AppendRow([]interface{}{
							o.Order.OrderID,
							o.Order.Price,
							o.Order.ExecutedQuantity,
							o.Order.CummulativeQuoteQuantity,
							o.Order.Side,
							o.Order.Status,
							time.Unix(o.Order.Time/1000, 0).Format(time.Stamp),
						}, table.RowConfig{
							AutoMerge: false,
						},
						)
						i++
					}
					//t.AppendFooter(table.Row{"", "", "", "", "", "", "", "", "", fmt.Sprintf("%11.8f", cumProfit), "", ""})
					if t.Length()>0 {
						t.Render()
					}
				}
				if dataChanged {
					binaClient.Store.Save()
				}
				if repeatDuration==0 {
					break
				}
				time.Sleep(repeatDuration)
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
				Name:        "telegrambotkey",
				Aliases:     []string{"tbk"},
				Usage:       "the bot key for telegram",
				EnvVars:     []string{"TELEGRAMBOTKEY"},
				FilePath:    "",
				Required:    false,
				Hidden:      false,
				TakesFile:   false,
				Value:       "",
				Destination: nil,
				HasBeenSet:  false,
			},
			&cli.StringFlag{
				Name:        "repeat",
				Aliases:     []string{"r"},
				Usage:       "run forever and check orders every <repeat> duration, eg. 30s or 1m or 1h, should somehow match the last duration ",
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
				Name:        "telegramchatid",
				Aliases:     []string{"tcid"},
				Usage:       "the chat id for telegram",
				EnvVars:     []string{"TELEGRAMCHATID"},
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
