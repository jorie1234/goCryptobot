package main

import (
	"github.com/adshao/go-binance/v2"
	"time"
)


type BinanceOrder struct {
	Order     binance.Order
	Relations []Relation
	FilledStatusSendToTelegram bool
}

func (bo BinanceOrder) OrderTimeYoungerThan(d time.Duration)  bool {
	if time.Unix(bo.Order.Time/1000, 0).Before(time.Now().Add(-d)) && time.Unix(bo.Order.UpdateTime/1000, 0).Before(time.Now().Add(-d))  {
		return false
	}
	return true
}