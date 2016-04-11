package main 

import (
    "encoding/json"
    "net/http"
    "fmt"
    "time"
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

func schedule(f func(), delaySeconds time.Duration) chan struct{} {
    ticker := time.NewTicker(delaySeconds * time.Second)
    
    quit := make(chan struct{})
    
    go func() {
        for {
            select {
                case <- ticker.C:
                    f()

                case <- quit:
                    ticker.Stop()
                    return
            }
        }
    }()

    return quit
}

type Minute struct {
    high float64
    low float64
    open float64
    close float64
    volume int
    ticker string
}

type MinuteHash struct {
    hash map[string]*Minute
}

func NewMinuteHash(tickers []string) *MinuteHash {
    hash := make(map[string]*Minute)

    for _, ticker := range tickers {
        hash[ticker] = &Minute{ticker: ticker}
    }

    return &MinuteHash{
        hash: hash,
    }
}

func (m *MinuteHash) add(t Trade) {
    minute := m.hash[t.Ticker]
    
    if minute.volume == 0 {
        minute.open = t.Price
    }

    if minute.low > t.Price {
        minute.low = t.Price 
    }

    if minute.high < t.Price {
        minute.high = t.Price
    }

    minute.volume = minute.volume + t.Shares
    minute.close = t.Price
}

func (m *MinuteHash) persistAndPublish(){
    minutes := make([]*Minute, 0)

    for ticker, tickMinute := range m.hash {
        minutes = append(minutes, tickMinute)
        m.publish(tickMinute)
        // reset hash
        m.hash[ticker] = &Minute{ticker: ticker}
    }

    m.persist(minutes)
}

func (m *MinuteHash) publish(tickMinute *Minute){

}

func (m *MinuteHash) persist(l []*Minute){
    fmt.Println("TICKER_SERVICE: ", l)
}

func main() {
    tickers := []string{"STOCK"}
    minuteHash := NewMinuteHash(tickers)

    schedule(minuteHash.persistAndPublish, 60)

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

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
        var payload [2]Trade
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&payload)
        if err != nil {
            panic(err)
        }

        minuteHash.add(payload[0])

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Status 200"))
    })

    http.ListenAndServe(":8003", nil)
}