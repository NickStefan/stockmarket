package main 

import (
    "encoding/json"
    "net/http"
)

type Trade struct {
    Actor string `json:"actor"`
    Shares int `json:"shares"`
    Ticker string `json:"ticker"`
    Price float64 `json:"price"`
    Intent string `json:"intent"`
    Kind string `json:"kind"`
    State  string `json:"state"`
    Time float64 `json:"time"`
}

func main() {

    hash := make(map[string][]Trade)

    // HOW SHOULD WE STORE THE DATA?
    // if speed not an issue, we would always add to mongo, and then aggregate the time series from mongo 
    // if speed is issue, we should store the series in redis, and then transform to requested series
    // we want speed, but dont want to store ALL of history in redis

    // need to impliment an aggregation framework on top of our redis store?
    // any data requested older than today, we then go to mongo?

    // older than n (e.g. today), use mongo aggregation query
    // today, get the pre-formatted ticks from redis

    // when we insert trades to redis ticker, preformat it for multiple ticker styles on multiple keys
    // then, end of day, do a job to put ticker data into actual DB

    // http://blog.mongodb.org/post/65517193370/schema-design-for-time-series-data-in-mongodb
    
    // {
    //       timestamp_hour: ISODate("2013-10-10T23:00:00.000Z"),
    //       type: “memory_used”,
    //       values: {
    //         0: { 0: {vol: 1000, price: 10.05}, …, 59: {vol: 1000, price: 10.05} },
    //         1: { 0: {vol: 1000, price: 10.05}, …, 59: {vol: 1000, price: 10.05} },
    //         …,
    //         58: { 0: {vol: 1000, price: 10.05}, …, 59: {vol: 1000, price: 10.05} },
    //         59: { 0: {vol: 1000, price: 10.05}, …, 59: {vol: 1000, price: 10.05} }
    //     }
    // }

    // {$set: {“values.59.59”: {vol: 1000, price: 10.05} } }

    // OR only store minutes (we dont need 1 sec or 10 sec charts)

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
    // map minute based data to correct format when querying server DB
    
    // real time chart:
    // then publish each new complete minute doc
    // front end will merge that into the data format of the queried for chart

    // 1 second ticker data will be separate channel, not persisted in a collection

    // on price data
    // throttle to every 1 second
    //      publish to ticker channel
    //      call chartData function, which is in turn throttled
    //          to update minute doc only when its been a minute(?)

    // question what to do when no trades are happening? update docs? 
    // if we could have "defaults" when aggregating docs OR saving docs, could help?

    // do we need an interval method to run every 60 seconds: http://stackoverflow.com/questions/16466320/is-there-a-way-to-do-repetitive-tasks-at-intervals-in-golang
    
    // every 60 seconds,
    //      for each stock symbols
    //          calc high, open, low, close, vol for last 60 seconds (cache last state to solve missing values?)
    //          clear cache?

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
        var payload [2]Trade
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&payload)
        if err != nil {
            panic(err)
        }

        ticker := payload[0].Ticker
        hash[ticker] = append(hash[ticker], payload[0])

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Status 200"))
    })   
}