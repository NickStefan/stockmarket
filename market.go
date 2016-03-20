package main

import (
	"fmt"
	"github.com/nickstefan/market/heap"
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
	state string // OPEN || FILLED || CANCELED
	timecreated string
	timeclosed string
}

func (b BaseOrder) lookup() string {
	return b.actor + b.timecreated
}

func (b BaseOrder) getOrder() BaseOrder {
	return b
}

type BuyLimit struct {
	bid float64
	BaseOrder
}

type SellLimit struct {
	ask float64
	BaseOrder
}

type BuyMarket struct {
	BaseOrder
}

type SellMarket struct {
	BaseOrder
}

// create a consistent interface for the different types of orders

type Order interface {
	price() float64
	lookup() string
	getOrder() BaseOrder
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
	buyQueue heap.Heap
	sellQueue heap.Heap
	buyHash map[string]*Order
	sellHash map[string]*Order
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		buyHash: make(map[string]*Order),
		sellHash: make(map[string]*Order),
		buyQueue: heap.Heap{Priority: "max"},
		sellQueue: heap.Heap{Priority: "min"},
	}
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
	for o.buyQueue.Peek().Value >= o.sellQueue.Peek().Value {
		o.buyQueue.Dequeue()
		o.sellQueue.Dequeue()
	}
}

// whats an algorithm to match buyers with sellers? simple case just using market orders

// orderBook.add( *order ) 
// puts the order into the buy or sell sub data structures
// then the order book attempts to fill that order
// if can fulfill
//    fills order, and removes the buy and sell orders that were affected by the fill
// 
// repeat on every orderBook.add( *order )

// orderBook.remove( *order )
// removes order from proper sub data structure

// so what does it mean to fill an order, and how could two data structures help us do that?

// an overlap between two heaps prioritized by price?
// but do we need constant time access to things deep inside the heap,
// due to order remove()?
// no, we mainly just need to pull things from the top of the buy heap and sell heap

// filling an order could be as simple as peeking at the buy heap and sell heap
// and asking if buy heap limit >= sell heap limit

// what about market orders? how would they prioritize into the heaps?
// do we need 4 heaps? or would market orders just always have priority?

// I think two heaps with market having highest priority works

// how would a limit order ever get filled if market is always higher priority?
// say the highest buy order is for 1000 shares limit $10.00,
// say the highest sell order is a market for 100 shares
// say the next highest sell order is a limit for 600 at $9.98.

// the trade would get partially filled at 100 shares @ $10.00 and
// 600 shares @ 9.99

// we can avoid having to traverse the heap to cancel things
// because were going to have a lookup hash where we store actual orders
// constant time cancelation, if the top of the heap gets looked up as canceled
// we remove it only then

func main() {
	orderBook := NewOrderBook()

	anOrder := SellLimit{
		ask: 10.05, 
		BaseOrder:BaseOrder{
			actor: "Tim", timecreated: "9:58am", intent: "SELL",
			shares: 100, state: "OPEN",
		},
	}

	anotherOrder := BuyLimit{
		bid: 10.00, 
		BaseOrder:BaseOrder{
			actor: "Bob", timecreated: "10:00am", intent: "BUY",
			shares: 100, state: "OPEN",
		},
	}
	orderBook.add(anOrder)
	orderBook.add(anotherOrder)

	fmt.Println(orderBook.buyHash)
	fmt.Println(orderBook.sellHash)
}