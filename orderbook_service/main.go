package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"time"
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

type Payload struct {
	Uuid   int      `json:"uuid"`
	Ticker string   `json:"ticker"`
	Orders []*Order `json:"orders"`
}

func main() {

	tickerUrl := "http://ticker:8080/ticker/trade"
	ledgerUrl := "http://ledger:8080/ledger/fill"

	redisAddress := "redis:6379"
	maxConnections := 10

	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", redisAddress)
		if err != nil {
			return nil, err
		}
		return c, err
	}, maxConnections)

	defer redisPool.Close()

	orderBook := NewOrderBook(redisPool)

	orderBook.setTradeHandler(func(t Trade, o Trade) {
		trade, err := json.Marshal([2]Trade{t, o})
		if err != nil {
			fmt.Println("orderbook_service: trade serialize http ", err)
		}

		go func() {
			ledgerResp, err := http.Post(ledgerUrl, "application/json", bytes.NewBuffer(trade))
			if err != nil {
				fmt.Println("orderbook_service: trade handler http ", err)
			}
			defer ledgerResp.Body.Close()
		}()

		go func() {
			tickerResp, err := http.Post(tickerUrl, "application/json", bytes.NewBuffer(trade))
			if err != nil {
				fmt.Println("orderbook_service: trade handler http ", err)
			}
			defer tickerResp.Body.Close()
		}()
	})

	makeTimeStamp := func() int64 {
		return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	}

	http.HandleFunc("/orderbook", func(w http.ResponseWriter, r *http.Request) {
		var payload Payload

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&payload)
		if err != nil {
			fmt.Println("orderbook_service: order handler http ", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Status 400"))
			return
		}

		if payload.Uuid == 500 || payload.Uuid == 1 {
			fmt.Println("receiving order", payload.Uuid, makeTimeStamp())
		}
		err = orderBook.Add(payload)
		if err != nil {
			fmt.Println("orderbook_service: orderbook Add  ", err)
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("Status 408"))
			return
		}
		if payload.Uuid == 500 || payload.Uuid == 1 {
			fmt.Println("order taken ", payload.Uuid, makeTimeStamp())
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})

	http.ListenAndServe(":8080", nil)
}
