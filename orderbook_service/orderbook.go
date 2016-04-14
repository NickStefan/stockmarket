package main

import (
	// "time"
	"github.com/nickstefan/market/orderbook_service/heap"
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
)

// how does a stock market organize the orders? Depth of Market or OrderBook

type OrderBook struct {
	lastPrice float64
	handleTrade tradehandler
	buyQueue heap.Heap
	sellQueue heap.Heap
	buyHash map[string]*Order
	sellHash map[string]*Order
}

type tradehandler func(Trade, Trade)

func NewOrderBook() *OrderBook {
	return &OrderBook{
		handleTrade: func(t Trade, o Trade) { },
		buyHash: make(map[string]*Order),
		sellHash: make(map[string]*Order),
		buyQueue: heap.Heap{Priority: "max"},
		sellQueue: heap.Heap{Priority: "min"},
	}
}

func (o *OrderBook) setTradeHandler(execTrade tradehandler) {
	o.handleTrade = execTrade
}

func (o *OrderBook) add(order *Order) {

	if order.Intent == "BUY" {
		o.buyHash[order.lookup()] = order
		o.buyQueue.Enqueue(&heap.Node{
			Value: order.price(),
			Lookup: order.lookup(),
		})

	} else if order.Intent == "SELL" {
		o.sellHash[order.lookup()] = order
		o.sellQueue.Enqueue(&heap.Node{
			Value: order.price(),
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

	} else if bKind == "LIMIT" && sKind == "LIMIT"{
		o.lastPrice = s.price() 
	
	} // else if both market, use last price

	return o.lastPrice
}

func (o *OrderBook) run() {
	buyTop := o.buyQueue.Peek()
	sellTop := o.sellQueue.Peek()
	
	for (buyTop != nil && sellTop != nil && buyTop.Value >= sellTop.Value) {
		
		buy := o.buyHash[ buyTop.Lookup ]
		sell := o.sellHash[ sellTop.Lookup ]
		price := o.negotiatePrice(buy, sell)

		if buy.Shares == sell.Shares {
			o.buyQueue.Dequeue()
			o.sellQueue.Dequeue()

			o.handleTrade(buy.fill(price), sell.fill(price))
			
			delete(o.buyHash, buyTop.Lookup)
			delete(o.sellHash, sellTop.Lookup)

		} else if buy.Shares < sell.Shares {
			o.buyQueue.Dequeue()
			remainderSell := sell.Shares - buy.Shares
			
			o.handleTrade(buy.fill(price), sell.partialFill(price, remainderSell))
 
			delete(o.buyHash, buyTop.Lookup)
		
		} else if buy.Shares > sell.Shares {
			o.sellQueue.Dequeue()
			remainderBuy := buy.Shares - sell.Shares
			
			o.handleTrade(sell.fill(price), sell.partialFill(price, remainderBuy))

			delete(o.sellHash, sellTop.Lookup)
		}
		
		buyTop = o.buyQueue.Peek()
		sellTop = o.sellQueue.Peek()
	}
}

func (o *OrderBook) addAll(payload Payload) {
	for _, order := range payload.Orders {
		o.add(order)
	}
}	

type Trade struct {
	Actor string `json:"actor"`
	Shares int `json:"shares"`
	Ticker string `json:"ticker"`
	Price float64 `json:"price"`
	Intent string `json:"intent"`
	Kind string `json:"kind"`
	State  string `json:"state"`
	Time int64 `json:"time"`
}

type Payload struct {
	Orders []*Order `json:"orders"`
}

func main() {

	ledgerUrl := "http://localhost:8002/fill"
	tickerUrl := "http://localhost:8003/"
	
	orderBook := NewOrderBook()

	orderBook.setTradeHandler(func (t Trade, o Trade) {
		trade, err := json.Marshal([2]Trade{t,o})
		if err != nil {
			panic(err)
		}

		_, err = http.Post(ledgerUrl, "application/json", bytes.NewBuffer(trade))
		if err != nil {
			panic(err)
		}

		_, err = http.Post(tickerUrl, "application/json", bytes.NewBuffer(trade))
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var payload Payload
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&payload)
		if err != nil {
			fmt.Println("ERR: ORDERBOOK_SERVICE")
			panic(err)
		}

		orderBook.addAll(payload)
		orderBook.run()

		w.WriteHeader(http.StatusOK)
        w.Write([]byte("Status 200"))
	})
	http.ListenAndServe(":8001", nil)
}