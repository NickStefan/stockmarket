package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
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

func schedule(f func() error, delaySeconds time.Duration) chan struct{} {
	ticker := time.NewTicker(delaySeconds * time.Second)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				f()

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

func main() {
	messageUrl := "http://web:8080/msg/ticker/"

	redisAddress := "redis:6379"
	maxConnections := 10

	mongoAddress := "192.168.99.100:27017"

	mongoSession, err := mgo.Dial(mongoAddress)
	if err != nil {
		fmt.Println("ticker_service: mongodb ", err)
	}
	//err = mongoSession.DB("tickerdb").DropDatabase()
	err = mongoSession.DB("tickerdb").C("ticks").EnsureIndex(mgo.Index{
		Key:        []string{"ticker", "date"},
		Background: true,
	})
	if err != nil {
		fmt.Println("ticker_service: mongodb ", err)
	}
	defer mongoSession.Close()

	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", redisAddress)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return c, err
	}, maxConnections)

	defer redisPool.Close()

	publisher := func(tickPeriod *Period) {
		tick, err := json.Marshal(struct {
			Payload *Period `json:"payload"`
			Api     string  `json:"api"`
			Version string  `json:"version"`
		}{
			Payload: tickPeriod,
			Api:     "ticker",
			Version: "1",
		})
		if err != nil {
			fmt.Println("ticker_service: publisher ", err)
		}

		messageResp, err := http.Post(messageUrl+tickPeriod.Ticker, "application/json", bytes.NewBuffer(tick))
		if err != nil {
			fmt.Println("ticker_service: publisher ", err)
			return
		}
		defer messageResp.Body.Close()
	}

	tickers := []string{"STOCK"}

	minuteRedis := NewPeriodHash(redisPool, "minute")
	minuteManager := NewPeriodManager(redisPool, minuteRedis, "minute")
	minuteManager.setTickers(tickers)
	minuteManager.initPeriods()
	minuteManager.setDB(mongoSession.DB("tickerdb"))

	secondRedis := NewPeriodHash(redisPool, "second")
	secondManager := NewPeriodManager(redisPool, secondRedis, "second")
	secondManager.setTickers(tickers)
	secondManager.initPeriods()
	secondManager.setPublisher(publisher)

	// TODO do this only if this server elected leader
	// if leader
	schedule(minuteManager.Persist, 60)
	schedule(secondManager.Publish, 1)
	// end if leader

	// TODO if unelect as leader
	// remove schedules
	// schedules return a chan that can be sent a quit message

	http.HandleFunc("/ticker/trade", func(w http.ResponseWriter, r *http.Request) {
		var payload [2]Trade
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&payload)
		if err != nil {
			fmt.Println("ticker_service: handle trade ", err)
		}

		// error handling?
		go minuteManager.add(payload[0])
		go secondManager.add(payload[0])

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})

	tickAggregator := TickAggregator{}

	tickAggregator.setDB(mongoSession.DB("tickerdb"))
	tickAggregator.setKV(minuteRedis)

	http.HandleFunc("/ticker/query", func(w http.ResponseWriter, r *http.Request) {
		var query Query
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		err := decoder.Decode(&query)
		if err != nil {
			fmt.Println("ticker_service: query handler", err)
		}

		results, err := tickAggregator.query(query)
		resultsJSON, err := json.Marshal(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Status", "200")
		w.Write(resultsJSON)
	})

	http.ListenAndServe(":8080", nil)
}
