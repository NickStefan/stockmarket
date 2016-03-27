package main

import (
	"fmt"
	"strconv"
	"time"
	"github.com/nickstefan/market/market_service/heap"
	"bytes"
	"encoding/json"
	"net/http"
)

// what makes up a stock market? competing buy and sell orders

// market buy order "i will buy 100 shares at the lowest available sell price"
// limit buy order "i will buy 100 shares at the lowest available price below ____"

// market sell order "i will sell 100 shares at the highest available price"
// limit sell order "i will sell 100 shares highest available price above _____"

type BaseOrder struct {
	shares int
	actor string // BOB
	intent string // BUY || SELL
	kind string // MARKET || LIMIT
	state string // OPEN || FILLED || CANCELED
	timecreated int64 // unix time
}

func (b *BaseOrder) lookup() string {
	return b.actor + strconv.FormatInt(b.timecreated, 10)
}

func (b *BaseOrder) getOrder() *BaseOrder {
	return b
}

func (b *BaseOrder) partialFill(price float64, newShares int) Trade {
	b.shares = newShares
	return Trade{Actor: b.actor, Shares: b.shares - newShares, Price: price, Intent: b.intent }
}

func (b *BaseOrder) fill(price float64) Trade {
	return Trade{Actor: b.actor, Shares: b.shares, Price: price, Intent: b.intent}
}

type BuyLimit struct {
	bid float64
	*BaseOrder
}

type SellLimit struct {
	ask float64
	*BaseOrder
}

type BuyMarket struct {
	*BaseOrder
}

type SellMarket struct {
	*BaseOrder
}

// create a consistent interface for the different types of orders

type Order interface {
	price() float64
	lookup() string
	getOrder() *BaseOrder
	partialFill(price float64, newShares int) Trade
	fill(price float64) Trade
}

func (b BuyLimit) price() float64 {
	return b.bid
}

func (s SellLimit) price() float64 {
	return s.ask
}

func (b BuyMarket) price() float64 {
	return 1000000.00
}

func (s SellMarket) price() float64 {
	return 0.00
}

// how does a stock market organize the orders? Depth of Market or OrderBook

type OrderBook struct {
	handleTrade tradehandler
	buyQueue heap.Heap
	sellQueue heap.Heap
	buyHash map[string]*Order
	sellHash map[string]*Order
}

type tradehandler func(Trade)

func NewOrderBook() *OrderBook {
	return &OrderBook{
		handleTrade: func(t Trade) { },
		buyHash: make(map[string]*Order),
		sellHash: make(map[string]*Order),
		buyQueue: heap.Heap{Priority: "max"},
		sellQueue: heap.Heap{Priority: "min"},
	}
}

func (o *OrderBook) setTradeHandler(execTrade tradehandler) {
	o.handleTrade = execTrade
}

func (o *OrderBook) add(order Order) {

	if order.getOrder().intent == "BUY" {
		o.buyHash[order.lookup()] = &order
		o.buyQueue.Enqueue(&heap.Node{
			Value: order.price(),
			Lookup: order.lookup(),
		})

	} else if order.getOrder().intent == "SELL" {
		o.sellHash[order.lookup()] = &order
		o.sellQueue.Enqueue(&heap.Node{
			Value: order.price(),
			Lookup: order.lookup(),
		})
	}
}

func (o *OrderBook) run() {
	buyTop := o.buyQueue.Peek()
	sellTop := o.sellQueue.Peek()
	
	for (buyTop != nil && sellTop != nil && buyTop.Value >= sellTop.Value) {
		
		buy := *(o.buyHash[ buyTop.Lookup ])
		sell := *(o.sellHash[ sellTop.Lookup ])

		if buy.getOrder().shares == sell.getOrder().shares {
			o.buyQueue.Dequeue()
			o.sellQueue.Dequeue()

			price := buy.price()
			o.handleTrade(buy.fill(price))
			o.handleTrade(sell.fill(sell.price()))
			
			delete(o.buyHash, buyTop.Lookup)
			delete(o.sellHash, sellTop.Lookup)

		} else if buy.getOrder().shares < sell.getOrder().shares {
			o.buyQueue.Dequeue()
			remainderSell := sell.getOrder().shares - buy.getOrder().shares
			
			price := buy.price()
			o.handleTrade(buy.fill(price))
			o.handleTrade(sell.partialFill(sell.price(), remainderSell))
 
			delete(o.buyHash, buyTop.Lookup)
		
		} else if buy.getOrder().shares > sell.getOrder().shares {
			o.sellQueue.Dequeue()
			remainderBuy := buy.getOrder().shares - sell.getOrder().shares
			
			price := buy.price()
			o.handleTrade(sell.fill(sell.price()))
			o.handleTrade(buy.partialFill(price, remainderBuy))

			delete(o.sellHash, sellTop.Lookup)
		}
		
		buyTop = o.buyQueue.Peek()
		sellTop = o.sellQueue.Peek()
	}
}


type Trade struct {
	Actor string
	Shares int
	Price float64
	Intent string
	Kind string
	State  string
}

func main() {

	orderBook := NewOrderBook()

	orderBook.setTradeHandler(func (t Trade) {
		fmt.Println("hello handler")
		url := "http://localhost:8000"
		trade, err := json.Marshal(t)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(trade))
		if err != nil {
			panic(err)
		}
		fmt.Println("response Status:", resp.Status)
	})

	anOrder := SellLimit{
		ask: 10.05, 
		BaseOrder: &BaseOrder{
			actor: "Tim", timecreated: time.Now().Unix(),
			intent: "SELL", shares: 100, state: "OPEN",
		},
	}

	anotherOrder := BuyLimit{
		bid: 10.10, 
		BaseOrder: &BaseOrder{
			actor: "Bob", timecreated: time.Now().Unix(),
			intent: "BUY", shares: 100, state: "OPEN",
		},
	}
	orderBook.add(anOrder)
	orderBook.add(anotherOrder)
	orderBook.run()
}