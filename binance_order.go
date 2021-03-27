package main

import (
	"bytes"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/jedib0t/go-pretty/v6/table"
)

type BinanceOrder struct {
	Order                      binance.Order
	Relations                  []Relation
	FilledStatusSendToTelegram bool
}

func (bo BinanceOrder) OrderTimeYoungerThan(d time.Duration) bool {
	if time.Unix(bo.Order.Time/1000, 0).Before(time.Now().Add(-d)) && time.Unix(bo.Order.UpdateTime/1000, 0).Before(time.Now().Add(-d)) {
		return false
	}
	return true
}

func (bo BinanceOrder) String() string {
	var buf bytes.Buffer
	t := table.NewWriter()
	t.SetOutputMirror(&buf)

	t.AppendSeparator()
	t.AppendHeader(table.Row{
		"Symbol",
		"OrderID",
		"Price",
		"Qnt",
		"ExQnt",
		"CumQuoteQnt",
		"Side",
		"Status",
		"Time",
	})
	t.SetStyle(table.StyleBold)

	t.AppendRow([]interface{}{
		bo.Order.Symbol,
		bo.Order.OrderID,
		bo.Order.Price,
		bo.Order.OrigQuantity,
		bo.Order.ExecutedQuantity,
		bo.Order.CummulativeQuoteQuantity,
		bo.Order.Side,
		bo.Order.Status,
		time.Unix(bo.Order.UpdateTime/1000, 0).Format(time.Stamp),
	}, table.RowConfig{
		AutoMerge: false,
	},
	)
	t.SetStyle(table.StyleBold)
	t.Render()
	return buf.String()
}
