package main

import (
	"github.com/nickstefan/market/orderbook_service/heap"
)

// how does a stock market organize the orders? Depth of Market or OrderBook

type OrderBook struct {
	lastPrice   float64
	handleTrade tradehandler
	buyQueue    heap.Heap
	sellQueue   heap.Heap
	buyHash     *OrderHash
	sellHash    *OrderHash
}

type tradehandler func(Trade, Trade)

func NewOrderBook() *OrderBook {
	return &OrderBook{
		handleTrade: func(t Trade, o Trade) {},
		buyHash:     NewOrderHash(),
		sellHash:    NewOrderHash(),
		buyQueue:    heap.Heap{Priority: "max"},
		sellQueue:   heap.Heap{Priority: "min"},
	}
}

func (o *OrderBook) setTradeHandler(execTrade tradehandler) {
	o.handleTrade = execTrade
}

func (o *OrderBook) add(order *Order) {
	if order.Intent == "BUY" {
		o.buyHash.set(order.lookup(), order)
		o.buyQueue.Enqueue(&heap.Node{
			Value:  order.price(),
			Lookup: order.lookup(),
		})

	} else if order.Intent == "SELL" {
		o.sellHash.set(order.lookup(), order)
		o.sellQueue.Enqueue(&heap.Node{
			Value:  order.price(),
			Lookup: order.lookup(),
		})
	}
}

func (o *OrderBook) negotiatePrice(b *Order, s *Order) float64 {
	bKind := b.Kind
	sKind := s.Kind

	if bKind == "MARKET" && sKind == "LIMIT" {
		o.lastPrice = s.price()

	} else if sKind == "MARKET" && bKind == "LIMIT" {
		o.lastPrice = b.price()

	} else if bKind == "LIMIT" && sKind == "LIMIT" {
		o.lastPrice = s.price()

	} // else if both market, use last price

	return o.lastPrice
}

func (o *OrderBook) run() {
	buyTop := o.buyQueue.Peek()
	sellTop := o.sellQueue.Peek()

	for buyTop != nil && sellTop != nil && buyTop.Value >= sellTop.Value {

		buy := o.buyHash.get(buyTop.Lookup)
		sell := o.sellHash.get(sellTop.Lookup)
		price := o.negotiatePrice(buy, sell)

		if buy.Shares == sell.Shares {
			o.buyQueue.Dequeue()
			o.sellQueue.Dequeue()

			o.handleTrade(buy.fill(price), sell.fill(price))

			o.buyHash.remove(buyTop.Lookup)
			o.sellHash.remove(sellTop.Lookup)

		} else if buy.Shares < sell.Shares {
			o.buyQueue.Dequeue()
			remainderSell := sell.Shares - buy.Shares

			o.handleTrade(buy.fill(price), sell.partialFill(price, remainderSell))

			o.buyHash.remove(buyTop.Lookup)

		} else if buy.Shares > sell.Shares {
			o.sellQueue.Dequeue()
			remainderBuy := buy.Shares - sell.Shares

			o.handleTrade(sell.fill(price), sell.partialFill(price, remainderBuy))

			o.sellHash.remove(sellTop.Lookup)
		}

		buyTop = o.buyQueue.Peek()
		sellTop = o.sellQueue.Peek()
	}
}
