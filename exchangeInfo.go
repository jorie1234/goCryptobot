package main

import (
	"fmt"
	"strings"

	"github.com/adshao/go-binance/v2"

	"github.com/urfave/cli/v2"
)

func InitExchangeInfo() *cli.Command {

	return &cli.Command{
		Name:         "exchangeinfo",
		Aliases:      []string{"ei"},
		Usage:        "prints exchange info",
		UsageText:    "",
		Description:  "",
		ArgsUsage:    "",
		Category:     "",
		BashComplete: nil,
		Before:       nil,
		After:        nil,
		Action: func(c *cli.Context) error {
			binaClient := GetBinanceClient(c)
			err := binaClient.GetExchangeInfo()
			if err != nil {
				fmt.Println("ExchangeInfo Error ", err)
				return nil
			}
			fmt.Println("successfully refreshed exchangeinfo")
			return nil
		},
		OnUsageError: nil,
		Subcommands:  nil,
	}
}

//func LoadExchangeInfo() *binance.ExchangeInfo {
//	var ei binance.ExchangeInfo
//	file := "exhangeinfo.json"
//	buffer, err := ioutil.ReadFile(file)
//	if err != nil {
//		return nil
//	}
//	err = json.Unmarshal(buffer, &ei)
//	if err != nil {
//		return nil
//	}
//	return nil
//}
//
//func GetSymbolFromExchangeInfo(ei *binance.ExchangeInfo, symbol string) *binance.Symbol {
//	for _, v := range ei.Symbols {
//		if v.Symbol == symbol {
//			return &v
//		}
//	}
//	return nil
//}

func GetLotSizeStepForSymbolFromExchangeInfo(ei *binance.ExchangeInfo, symbol string) string {
	var stepSize string
	var lotFound bool
	for _, v := range ei.Symbols {
		if v.Symbol == symbol {
			for _, f := range v.Filters {
				for k, ff := range f {
					if k == "stepSize" {
						stepSize = ff.(string)
						if lotFound {
							return stepSize
						}
					}
					if k == "filterType" {
						if ff == "LOT_SIZE" {
							lotFound = true
							if len(stepSize) > 0 {
								return stepSize
							}
						}
					}
					//fmt.Printf("%s %s\n", k, ff)
				}
			}
			return ""
		}
	}
	return ""
}

func TrimQuantityToLotSize(quantity, lotSize string) string {
	qa := strings.Split(quantity, ".")
	la := strings.Split(lotSize, ".")
	if len(qa) == 1 || len(qa[1]) == 0 {
		return qa[0]
	}
	result := qa[0] + "."
	for p, v := range la[1] {
		char := fmt.Sprintf("%c", v)
		//fmt.Printf("char is %s\n", char)
		result = result + fmt.Sprintf("%c", qa[1][p])
		if char == "1" || p >= len(qa[1])-1 {
			return result
		}
	}
	return ""
}
