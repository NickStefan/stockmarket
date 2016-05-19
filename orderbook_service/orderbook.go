package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/nickstefan/market/orderbook_service/heap"
)

// how does a stock market organize the orders? Depth of Market or OrderBook

type OrderBook struct {
	lastPrice   float64
	handleTrade tradehandler
	orderQueue  *OrderQueue
	orderHash   *OrderHash
}

type tradehandler func(Trade, Trade)

func NewOrderBook(pool *redis.Pool) *OrderBook {
	return &OrderBook{
		handleTrade: func(t Trade, o Trade) {},
		orderHash:   NewOrderHash(pool, ""),
		orderQueue:  NewOrderQueue(pool),
	}
}

func (o *OrderBook) setEnv(env string) {
	o.orderQueue.setEnv(env)
	o.orderHash.setEnv(env)
}

func (o *OrderBook) setTradeHandler(execTrade tradehandler) {
	o.handleTrade = execTrade
}

func (o *OrderBook) add(order *Order) {
	o.orderHash.set(order.lookup(), order)
	o.orderQueue.Enqueue(order.Intent+order.Ticker, &heap.Node{
		Value:  order.price(),
		Lookup: order.lookup(),
	})
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

func (o *OrderBook) run(ticker string) {
	buyTop := o.orderQueue.Peek("BUY" + ticker)
	sellTop := o.orderQueue.Peek("SELL" + ticker)

	for buyTop != nil && sellTop != nil && buyTop.Value >= sellTop.Value {

		buy := o.orderHash.get(buyTop.Lookup)
		sell := o.orderHash.get(sellTop.Lookup)
		price := o.negotiatePrice(buy, sell)

		if buy.Shares == sell.Shares {
			o.orderQueue.Dequeue("BUY" + ticker)
			o.orderQueue.Dequeue("SELL" + ticker)

			o.handleTrade(buy.fill(price), sell.fill(price))

			o.orderHash.remove(buyTop.Lookup)
			o.orderHash.remove(sellTop.Lookup)

		} else if buy.Shares < sell.Shares {
			o.orderQueue.Dequeue("BUY" + ticker)
			remainderSell := sell.Shares - buy.Shares

			o.handleTrade(buy.fill(price), sell.partialFill(price, remainderSell))

			o.orderHash.remove(buyTop.Lookup)

		} else if buy.Shares > sell.Shares {
			o.orderQueue.Dequeue("SELL" + ticker)
			remainderBuy := buy.Shares - sell.Shares

			o.handleTrade(sell.fill(price), sell.partialFill(price, remainderBuy))

			o.orderHash.remove(sellTop.Lookup)
		}

		buyTop = o.orderQueue.Peek("BUY" + ticker)
		sellTop = o.orderQueue.Peek("SELL" + ticker)
	}
}
