package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/nickstefan/market/orderbook_service/heap"
	"gopkg.in/redsync.v1"
	"sync"
	"time"
)

type OrderBook struct {
	lockMap       map[string]*Locker
	redsync       *redsync.Redsync
	lastPrice     float64
	handleTrade   tradehandler
	publishBidAsk bidaskpublisher
	orderQueue    *OrderQueue
	orderHash     *OrderHash
}

type tradehandler func(Trade, Trade)
type bidaskpublisher func(bid *Order, ask *Order)

func NewOrderBook(pool *redis.Pool) *OrderBook {
	return &OrderBook{
		lockMap:       make(map[string]*Locker),
		redsync:       redsync.New([]redsync.Pool{pool}),
		handleTrade:   func(t Trade, o Trade) {},
		publishBidAsk: func(bid *Order, ask *Order) {},
		orderHash:     NewOrderHash(pool, ""),
		orderQueue:    NewOrderQueue(pool),
	}
}

func (o *OrderBook) setEnv(env string) {
	o.orderQueue.setEnv(env)
	o.orderHash.setEnv(env)
}

func (o *OrderBook) setHandleTrade(execTrade tradehandler) {
	o.handleTrade = execTrade
}

func (o *OrderBook) setPublishBidAsk(publisher bidaskpublisher) {
	o.publishBidAsk = publisher
}

func (o *OrderBook) getLocker(ticker string) *Locker {
	if nil != o.lockMap[ticker] {
		return o.lockMap[ticker]
	} else {
		redLockMutex := o.redsync.NewMutex("orderbook_service" + ticker)
		redsync.SetRetryDelay(5 * time.Millisecond).Apply(redLockMutex)
		redsync.SetExpiry(500 * time.Millisecond).Apply(redLockMutex)
		redsync.SetTries(50).Apply(redLockMutex)

		o.lockMap[ticker] = &Locker{
			name:    "orderbook_service" + ticker,
			mutLock: &sync.Mutex{},
			redLock: redLockMutex,
		}
		return o.lockMap[ticker]
	}
}

func (o *OrderBook) Add(payload Payload) error {
	locker := o.getLocker(payload.Ticker)
	err := locker.Lock()
	if err != nil {
		return err
	}
	defer locker.Unlock()

	for _, order := range payload.Orders {
		o.add(order)
	}
	o.run(payload.Ticker)
	return nil
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

			go o.handleTrade(buy.fill(price), sell.fill(price))

			o.orderHash.remove(buyTop.Lookup)
			o.orderHash.remove(sellTop.Lookup)

		} else if buy.Shares < sell.Shares {
			o.orderQueue.Dequeue("BUY" + ticker)
			remainderSell := sell.Shares - buy.Shares

			go o.handleTrade(buy.fill(price), sell.partialFill(price, remainderSell))

			o.orderHash.remove(buyTop.Lookup)

		} else if buy.Shares > sell.Shares {
			o.orderQueue.Dequeue("SELL" + ticker)
			remainderBuy := buy.Shares - sell.Shares

			go o.handleTrade(sell.fill(price), sell.partialFill(price, remainderBuy))

			o.orderHash.remove(sellTop.Lookup)
		}

		buyTop = o.orderQueue.Peek("BUY" + ticker)
		sellTop = o.orderQueue.Peek("SELL" + ticker)
	}

	o.publishOrderBook(ticker, buyTop, sellTop)
}

func (o *OrderBook) publishOrderBook(ticker string, buyTop *heap.Node, sellTop *heap.Node) {
	bid := &Order{Ticker: ticker}
	ask := &Order{Ticker: ticker}

	if buyTop != nil {
		bid = o.orderHash.get(buyTop.Lookup)
	}
	if sellTop != nil {
		ask = o.orderHash.get(sellTop.Lookup)
	}
	go o.publishBidAsk(bid, ask)
}
