package main

import (
    // "fmt"
    "gopkg.in/mgo.v2"
    "time"
)

type Minute struct {
    High float64 `json:"high"`
    Low float64 `json:"low"`
    Open float64 `json:"open"`
    Close float64 `json:"close"`
    Volume int `json:"volume"`
    Ticker string `json:"ticker"`
    Time int64 `json:"time"`
}

type MinuteHash struct {
    hash map[string]*Minute
    db *mgo.Database
}

func NewMinuteHash(tickers []string) *MinuteHash {
    hash := make(map[string]*Minute)

    for _, ticker := range tickers {
        hash[ticker] = &Minute{Ticker: ticker}
    }

    return &MinuteHash{
        hash: hash,
    }
}

func (m *MinuteHash) setDB(db *mgo.Database) {
    m.db = db
}

func (m *MinuteHash) add(t Trade) {
    minute := m.hash[t.Ticker]
    
    if minute.Volume == 0 {
        minute.Open = t.Price
        minute.High = t.Price
        minute.Low = t.Price
        minute.Close = t.Price
        minute.Time = time.Now().Unix()
    }

    if minute.Low > t.Price {
        minute.Low = t.Price 
    }

    if minute.High < t.Price {
        minute.High = t.Price
    }

    minute.Volume = minute.Volume + t.Shares
    minute.Close = t.Price
}

func (m *MinuteHash) persistAndPublish(){
    minutes := make([]interface{}, 0)

    for ticker, tickMinute := range m.hash {
        minutes = append(minutes, tickMinute)
        m.publish(tickMinute)
        // reset hash
        m.hash[ticker] = &Minute{Ticker: ticker}
    }

    m.persist(minutes)
}

func (m *MinuteHash) publish(tickMinute *Minute){
    // fmt.Println("TICKER_SERVICE: ", tickMinute)
}

func (m *MinuteHash) persist(l []interface{}){
    // put it all into mongodb
    c := m.db.C("ticks")
    err := c.Insert(l...)
    if err != nil {
        panic(err)
    }
}