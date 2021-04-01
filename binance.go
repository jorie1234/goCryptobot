package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/patrickmn/go-cache"
)

type BinanceClient struct {
	client *binance.Client
	cache  *cache.Cache
	Store  BinanceOrderStore
}
type Relation struct {
	BuyOrderID  int64
	SellOrderID int64
	Type        int
	Percent     float32
	Profit      float32
	Mult        float64
}

type BinanceOrderStore struct {
	Orders []BinanceOrder
}

func (bos *BinanceOrderStore) Load() error {
	file := "binanceorders.json"
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buffer, &bos.Orders)
	if err != nil {
		return err
	}
	fmt.Printf("Loaded %d orders\n", len(bos.Orders))
	return nil
}

func (bos BinanceOrderStore) Save() error {
	file := "binanceorders.json"
	r, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()

	orderJson, err := json.MarshalIndent(bos.Orders, "", "  ")
	if err != nil {
		fmt.Printf("Cannot encode to JSON %s", err)
	}
	_, _ = fmt.Fprintf(r, "%s", orderJson)

	return nil
}

func NewBinanceClient(binaAPIKey, binaSecretKey string) *BinanceClient {

	bc := BinanceClient{
		client: binance.NewClient(binaAPIKey, binaSecretKey),
		cache:  cache.New(5*time.Minute, 10*time.Minute),
	}
	err := bc.Store.Load()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &bc
}

func (bc *BinanceClient) GetAveragePrice(symbol string) (*binance.AvgPrice, error) {

	avgPriceCached, found := bc.cache.Get(symbol)
	if found {
		return avgPriceCached.(*binance.AvgPrice), nil
	}

	avgPrice, err := bc.client.NewAveragePriceService().Symbol(symbol).Do(context.Background())

	if err != nil {
		//fmt.Printf("Cannot get price for Symbol %s : %v\n", symbol, err)
		return nil, err
	}
	bc.cache.Set(symbol, avgPrice, cache.DefaultExpiration)
	return avgPrice, err
}

func (bos BinanceOrderStore) GetOrderByID(id int64) *BinanceOrder {
	for _, v := range bos.Orders {
		//fmt.Printf("check order %d vs order %d\n", v.Order.OrderID, id)
		if v.Order.OrderID == id {
			return &v
		}
	}
	return nil
}

func (bos BinanceOrderStore) SetSellForOrder(buyOrder int64, sellOrder int64, mult float64) {
	for k, v := range bos.Orders {
		if v.Order.OrderID == buyOrder {
			bos.Orders[k].Relations = append(bos.Orders[k].Relations, Relation{
				BuyOrderID:  buyOrder,
				SellOrderID: sellOrder,
				Type:        0,
				Percent:     1.0,
				Profit:      0,
				Mult:        mult,
			})
			return
		}
	}
}
func (bc *BinanceClient) PostSellForOrder(orderID int64, symbol string, mult float64) (*binance.Order, error) {

	order := bc.Store.GetOrderByID(orderID)
	if order == nil {
		return nil, fmt.Errorf("cannot find order with id %d\n", orderID)
	}

	if bc.DoesASellOrderExistForThisOrder(order) {
		return nil, fmt.Errorf("order with id %d already sold\n", orderID)
	}
	if order.Order.Side != binance.SideTypeBuy {
		return nil, fmt.Errorf("order with id %d is of type %s \n", orderID, order.Order.Side)
	}
	if order.Order.Symbol != symbol {
		return nil, fmt.Errorf("order with id %d has symbol %s not symbol %s \n", orderID, order.Order.Symbol, symbol)
	}

	ord, err := bc.client.NewGetOrderService().OrderID(orderID).Symbol(symbol).
		Do(context.Background())

	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	//bc.cache.Set(s, avgPrice, cache.DefaultExpiration)
	//return avgPrice, err
	fmt.Printf("fetched order %+#v\n", ord)
	Price, _ := strconv.ParseFloat(ord.Price, 8)
	Price = Price * mult
	PriceString := fmt.Sprintf("%f", Price)
	quant, _ := strconv.ParseFloat(ord.ExecutedQuantity, 8)
	cumQuote, _ := strconv.ParseFloat(ord.CummulativeQuoteQuantity, 8)
	fmt.Printf("Sell orderID for Price %s quant %s sell will be at %.2f€ profit will be %.2f€\n", PriceString, ord.ExecutedQuantity, Price*quant, Price*quant-cumQuote)
	if !Confirm("Perform sell ?") {
		return nil, nil
	}
	orderResponse, err := bc.client.NewCreateOrderService().
		Symbol(symbol).
		Type(binance.OrderTypeLimit).
		Quantity(ord.ExecutedQuantity).
		Price(PriceString).
		Side(binance.SideTypeSell).
		TimeInForce(binance.TimeInForceTypeGTC).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	fmt.Printf("OrderResponse %+#v", orderResponse)
	bc.Store.SetSellForOrder(orderID, orderResponse.OrderID, mult)
	err = bc.Store.Save()
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	return ord, nil
}

func (bc *BinanceClient) ListOrders(symbol string) {
	orders, err := bc.client.NewListOrdersService().Symbol(symbol).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return //orders, nil
	}
	log.Printf("found %d orders for symbol %s\n", len(orders), symbol)
	for _, v := range orders {
		bc.InsertOrder(v)
	}
	err = bc.Store.Save()
	if err != nil {
		fmt.Println(err)
		return //orders, nil
	}

}

func (bc *BinanceClient) DoesASellOrderExistForThisOrder(order *BinanceOrder) bool {
	for _, r := range order.Relations {
		if r.SellOrderID > 1 {
			return true
		}
	}
	return false
}

func (bc *BinanceClient) GetRawSellOrderID(order *BinanceOrder) int64 {
	for _, r := range order.Relations {
		if r.SellOrderID > 0 {
			return r.SellOrderID
		}
	}
	return 0
}

func (bc *BinanceClient) IsRelationYoungerThan(o BinanceOrder, d time.Duration) bool {
	for _, r := range o.Relations {
		order := bc.GetOrderByID(r.SellOrderID)
		if order != nil {
			if order.OrderTimeYoungerThan(d) {
				return true
			}
		}
	}
	return false
}

func (bc *BinanceClient) GetSellOrderIDforOrder(order *BinanceOrder) int64 {
	for _, r := range order.Relations {
		if r.SellOrderID > 1 {
			return r.SellOrderID
		}
	}
	return 0
}
func (bc *BinanceClient) GetBuyOrderIDforOrder(order *BinanceOrder) *BinanceOrder {
	for _, r := range bc.Store.Orders {
		for _, relation := range r.Relations {
			if relation.SellOrderID == order.Order.OrderID {
				return &r
			}
		}
	}
	return nil
}

func (bc *BinanceClient) GetOrderByID(id int64) *BinanceOrder {
	for _, r := range bc.Store.Orders {
		if r.Order.OrderID == id {
			return &r
		}
	}
	return nil
}
func (bc *BinanceClient) InsertOrder(o *binance.Order) {
	found := false
	for k, v := range bc.Store.Orders {
		if v.Order.OrderID == o.OrderID {
			found = true
			bc.Store.Orders[k].Order = *o
			//fmt.Printf("found order %d\n", v.Order.OrderID)
			continue
		}
	}
	if !found {
		bo := BinanceOrder{
			Order: *o,
			Relations: []Relation{{
				BuyOrderID:  0,
				SellOrderID: 0,
				Type:        0,
				Percent:     0,
				Profit:      0,
			},
			}}
		bc.Store.Orders = append(bc.Store.Orders, bo)
	}
	//fmt.Printf("now %d orders in store", len(bc.Store.Orders))
}

func (bc *BinanceClient) GetStatusByOderID(id int64) binance.OrderStatusType {
	for _, v := range bc.Store.Orders {
		if v.Order.OrderID == id {
			return v.Order.Status
		}
	}
	return ""
}

func RemoveRelation(s []Relation, index int) []Relation {
	return append(s[:index], s[index+1:]...)
}

func (bc *BinanceClient) CancelOrder(id int64) error {
	bo := bc.GetOrderByID(id)
	if bo == nil {
		return fmt.Errorf("cannot find order with id %d", id)
	}
	res, err := bc.client.NewCancelOrderService().Symbol(bo.Order.Symbol).OrderID(bo.Order.OrderID).Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("CancelOrderResponse %+#v\n", res)
	buyOrder := bc.GetBuyOrderIDforOrder(bo)
	if buyOrder == nil {
		fmt.Printf("Order %d has no associated buy order", id)
		return nil
	}
	for k, v := range buyOrder.Relations {
		if v.SellOrderID == id {
			buyOrder.Relations = RemoveRelation(buyOrder.Relations, k)
			fmt.Printf("delete order %d from Buy order %d", id, buyOrder.Order.OrderID)
		}
	}
	return nil
}

func (bc *BinanceClient) ConnectSellOrderWithBuyOrder(sellorderID, buyorderID int64) error {
	buyorder := bc.GetOrderByID(buyorderID)
	if buyorder == nil {
		return fmt.Errorf("cannot find order for buyorderID %d", buyorderID)
	}
	sellorder := bc.GetOrderByID(sellorderID)
	if sellorder == nil {
		return fmt.Errorf("cannot find order for sellorderID %d", buyorderID)
	}
	if buyorder.Order.Side != binance.SideTypeBuy {
		return fmt.Errorf("order for buyorderID %d is not a buy order, order is %s", buyorderID, buyorder.Order.Side)
	}
	if sellorder.Order.Side != binance.SideTypeSell {
		return fmt.Errorf("order for sellorderID %d is not a sell order, order is %s", buyorderID, sellorder.Order.Side)
	}
	bc.Store.SetSellForOrder(buyorderID, sellorderID, 0.0)
	return nil
}
