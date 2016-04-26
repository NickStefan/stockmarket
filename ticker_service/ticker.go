package main

import (
	"encoding/json"
	"fmt"
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
	url := "mongodb://localhost"
	session, err := mgo.Dial(url)
	err = session.DB("tickerdb").DropDatabase()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	tickers := []string{"STOCK"}
	periodHash := NewPeriodHash(tickers)
	periodHash.setDB(session.DB("tickerdb"))

	schedule(periodHash.persistAndPublish, 60)

	// {
	//       timestamp_hour: ISODate("2013-10-10T23:00:00.000Z"),
	//       type: “memory_used”,
	//       values: {
	//         0: {vol: 1000, open: 10.05, close: 10.55, high: 11.00, low: 10.00 },
	//         1: {vol: 1000, open: 10.55, close: 10.60, high: 11.00, low: 10.50 },
	//         …,
	//         58: {vol: 1000, open: 10.65, close: 10.80, high: 11.00, low: 10.60 },
	//         59: {vol: 1000, open: 10.65, close: 11.55, high: 12.00, low: 10.00 }
	//       }
	//     }
	// }
	// {$set: {“values.59”: {vol: 1000, open: 10.65, close: 11.55, high: 12.00, low: 10.00 } }

	// basic charting:
	// aggregate minute based data to correct format time ranges

	// what will an aggregation query look like? "I want the last <X>(5?) <range>(days?) <duration>(hour?) chart"
	// $match stock <X> * <range>
	// $unwind values
	// $groupby <duration> (want array of duration groups)
	// $sum vol, $min low, $max hi, open, close of first?
	// $map appropriate time stamps?

	// real time chart:
	// front end will merge published minute info into the chart's range

	// 1 second ticker data will be separate channel, not persisted in a collection

	// on price data
	//      - publish to ticker channel, rate limit to one / second
	//      - add to cache of last 60 seconds of trades

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var payload [2]Trade
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&payload)
		if err != nil {
			fmt.Println("ERR: TICKER_SERVICE")
			panic(err)
		}

		periodHash.add(payload[0])

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})

	http.ListenAndServe(":8003", nil)
}
