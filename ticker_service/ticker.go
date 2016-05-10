package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
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

func schedule(f func(), delaySeconds time.Duration) chan struct{} {
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
	messageUrl := "http://127.0.0.1:8004/msg/ticker/"
	url := "mongodb://localhost"

	session, err := mgo.Dial(url)
	// db.temperature.ensureIndex({'date': 1 })'}
	//err = session.DB("tickerdb").DropDatabase()
	if err != nil {
		fmt.Println("TODO: ticker_service fault tolerance needed; ", err)
	}
	defer session.Close()

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
			fmt.Println("TODO: ticker_service fault tolerance needed; ", err)
		}

		messageResp, err := http.Post(messageUrl+tickPeriod.Ticker, "application/json", bytes.NewBuffer(tick))
		if err != nil {
			fmt.Println("TODO: ticker_service fault tolerance needed; ", err)
		}
		defer messageResp.Body.Close()
	}

	tickers := []string{"STOCK"}

	minuteHash := NewPeriodHash(tickers)
	minuteHash.setDB(session.DB("tickerdb"))
	schedule(minuteHash.Persist, 60)

	secondHash := NewPeriodHash(tickers)
	secondHash.setPublisher(publisher)
	schedule(secondHash.Publish, 1)

	http.HandleFunc("/trade", func(w http.ResponseWriter, r *http.Request) {
		var payload [2]Trade
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&payload)
		if err != nil {
			fmt.Println("TODO: ticker_service fault tolerance needed; ", err)
		}

		minuteHash.add(payload[0])
		secondHash.add(payload[0])

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})

	tickAggregator := TickAggregator{}
	tickAggregator.setDB(session.DB("tickerdb"))

	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		results := tickAggregator.query(Query{
			TickerName:   "STOCK",
			Periods:      2,
			PeriodNumber: 1,
			PeriodName:   "minute",
		})

		resultsJSON, err := json.Marshal(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Status", "200")
		w.Write(resultsJSON)
	})

	http.ListenAndServe(":8003", nil)
}
