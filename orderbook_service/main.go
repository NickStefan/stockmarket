package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/redsync.v1"
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

	redisAddress := ":6379"
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

	var mutex sync.Mutex

	redisLock := redsync.New([]redsync.Pool{redisPool})

	orderBook.setTradeHandler(func(t Trade, o Trade) {
		// fmt.Println("\n TRADE", t.Price, "\n")
		trade, err := json.Marshal([2]Trade{t, o})
		if err != nil {
			fmt.Println("orderbook_service: trade serialize http ", err)
		}

		ledgerResp, err := http.Post(ledgerUrl, "application/json", bytes.NewBuffer(trade))
		if err != nil {
			fmt.Println("orderbook_service: trade handler http ", err)
		}
		defer ledgerResp.Body.Close()

		tickerResp, err := http.Post(tickerUrl, "application/json", bytes.NewBuffer(trade))
		if err != nil {
			fmt.Println("orderbook_service: trade handler http ", err)
		}
		defer tickerResp.Body.Close()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Ticker string   `json:"ticker"`
			Orders []*Order `json:"orders"`
		}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&payload)
		if err != nil {
			fmt.Println("orderbook_service: order handler http ", err)
		}

		// LOCKING
		mutex.Lock()
		defer mutex.Unlock()
		// do we need BOTH the mutex AND the redlock???
		rMutex := redisLock.NewMutex(payload.Ticker)
		err = rMutex.Lock()
		if err != nil {
			panic(err)
		}
		defer rMutex.Unlock()
		// END LOCKING

		for _, order := range payload.Orders {
			orderBook.add(order)
		}
		orderBook.run(payload.Ticker)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})

	http.ListenAndServe(":8001", nil)
}
