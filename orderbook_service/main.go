package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Trade struct {
	Actor  string  `json:"actor"`
	Shares int     `json:"shares"`
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
	Intent string  `json:"intent"`
	Kind   string  `json:"kind"`
	State  string  `json:"state"`
	Time   int64   `json:"time"`
}

func main() {

	ledgerUrl := "http://127.0.0.1:8002/fill"
	tickerUrl := "http://127.0.0.1:8003/trade"

	orderBook := NewOrderBook()
	var mutex sync.Mutex

	orderBook.setTradeHandler(func(t Trade, o Trade) {
		// fmt.Println("\n TRADE", t.Price, "\n")
		trade, err := json.Marshal([2]Trade{t, o})
		if err != nil {
			fmt.Println("TODO: orderbook fault tolerance needed; ", err)
		}

		// TODO
		// TODO
		// TODO
		// TODO
		// TODO
		// should put some logs here to test if ledger Post
		// blocks the ticker post
		ledgerResp, err := http.Post(ledgerUrl, "application/json", bytes.NewBuffer(trade))
		if err != nil {
			fmt.Println("TODO: orderbook fault tolerance needed; ", err)
		}
		defer ledgerResp.Body.Close()

		tickerResp, err := http.Post(tickerUrl, "application/json", bytes.NewBuffer(trade))
		if err != nil {
			fmt.Println("TODO: orderbook fault tolerance needed; ", err)
		}
		defer tickerResp.Body.Close()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Orders []*Order `json:"orders"`
		}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&payload)
		if err != nil {
			fmt.Println("TODO: orderbook fault tolerance needed; ", err)
		}

		mutex.Lock()
		defer mutex.Unlock()
		for _, order := range payload.Orders {
			// fmt.Println(order.Intent, "ORDER", order.price())
			orderBook.add(order)
		}
		orderBook.run()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})

	http.ListenAndServe(":8001", nil)
}