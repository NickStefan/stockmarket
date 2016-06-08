package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"runtime"
	"time"
)

type Payload struct {
	Uuid   int      `json:"uuid"`
	Ticker string   `json:"ticker"`
	Orders []*Order `json:"orders"`
}

func makeTimeStamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func main() {

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

	messageUrl := "http://web:8080/msg/ticker/"
	tickerUrl := "http://ticker:8080/ticker/trade"
	//ledgerUrl := "http://ledger:8080/ledger/fill"

	orderBook.setHandleTrade(func(t Trade, o Trade) {
		//go func() {
		//trades, err := json.Marshal([2]Trade{t, o})
		//if err != nil {
		//fmt.Println("orderbook_service: ledger serialize http ", err)
		//}

		//ledgerResp, err := http.Post(ledgerUrl, "application/json", bytes.NewBuffer(trades))
		//if err != nil {
		//fmt.Println("orderbook_service: trade handler http ", err)
		//return
		//}
		//defer ledgerResp.Body.Close()
		//}()

		go func() {
			tick, err := json.Marshal(struct {
				Payload AnonymizedTrade `json:"payload"`
				Api     string          `json:"api"`
				Version string          `json:"version"`
			}{
				Payload: AnonymizeTrade(t),
				Api:     "ticker",
				Version: "1",
			})
			if err != nil {
				fmt.Println("orderbook_service: msg seriailze ", err)
			}

			if t.Price == 2 {
				fmt.Println("msg started", makeTimeStamp())
			}
			if t.Price == 70 {
				fmt.Println("msg ending", makeTimeStamp())
			}
			messageResp, err := http.Post(messageUrl+t.Ticker, "application/json", bytes.NewBuffer(tick))
			if err != nil {
				fmt.Println("orderbook_service: msg http ", err)
				return
			}
			defer messageResp.Body.Close()
		}()

		go func() {
			trade, err := json.Marshal(AnonymizeTrade(t))
			if err != nil {
				fmt.Println("orderbook_service: ticker serialize", err)
			}

			tickerResp, err := http.Post(tickerUrl, "application/json", bytes.NewBuffer(trade))
			if err != nil {
				fmt.Println("orderbook_service: ticker http ", err)
				return
			}
			defer tickerResp.Body.Close()
		}()
	})

	orderBook.setPublishBidAsk(func(bid *Order, ask *Order) {
		bidAsk, err := json.Marshal(struct {
			Payload []AnonymizedOrder `json:"payload"`
			Api     string            `json:"api"`
			Version string            `json:"version"`
		}{
			Payload: []AnonymizedOrder{
				AnonymizeOrder(bid),
				AnonymizeOrder(ask),
			},
			Api:     "bid-ask",
			Version: "1",
		})
		if err != nil {
			fmt.Println("orderbook_service: bid-ask seriailze ", err)
		}

		messageResp, err := http.Post(messageUrl+bid.Ticker, "application/json", bytes.NewBuffer(bidAsk))
		if err != nil {
			fmt.Println("orderbook_service: bid-ask http ", err)
			return
		}
		defer messageResp.Body.Close()
	})

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

	http.HandleFunc("/orderbook/cpu", func(w http.ResponseWriter, r *http.Request) {
		maxProcs := runtime.GOMAXPROCS(0)
		numCPU := runtime.NumCPU()
		fmt.Println("Procs and cpu ", maxProcs, numCPU)

		info, err := json.Marshal(struct {
			Cpu   int
			Procs int
		}{
			Cpu:   numCPU,
			Procs: maxProcs,
		})
		if err != nil {
			fmt.Println(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(info)
	})

	http.ListenAndServe(":8080", nil)
}
