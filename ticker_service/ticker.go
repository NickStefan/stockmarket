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