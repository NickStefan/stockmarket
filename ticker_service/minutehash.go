package main

import (
  "fmt"
)

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
        minute.high = t.Price
        minute.low = t.Price
        minute.close = t.Price
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