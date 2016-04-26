package main

import (
	// "fmt"
	"gopkg.in/mgo.v2"
	"sync"
	"time"
)

type Period struct {
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Volume int     `json:"volume"`
	Ticker string  `json:"ticker"`
	Time   int64   `json:"time"`
}

type PeriodHash struct {
	hash map[string]*Period
	db   *mgo.Database
	sync.RWMutex
}

func NewPeriodHash(tickers []string) *PeriodHash {
	hash := make(map[string]*Period)

	for _, ticker := range tickers {
		hash[ticker] = &Period{Ticker: ticker}
	}

	return &PeriodHash{
		hash: hash,
	}
}

func (m *PeriodHash) setDB(db *mgo.Database) {
	m.db = db
}

func (m *PeriodHash) add(t Trade) {
	m.Lock()
	defer m.Unlock()

	period := m.hash[t.Ticker]

	if period.Volume == 0 {
		period.Open = t.Price
		period.High = t.Price
		period.Low = t.Price
		period.Close = t.Price
		period.Time = time.Now().Unix()
	}

	if period.Low > t.Price {
		period.Low = t.Price
	}

	if period.High < t.Price {
		period.High = t.Price
	}

	period.Volume = period.Volume + t.Shares
	period.Close = t.Price
}

func (m *PeriodHash) persistAndPublish() {
	m.RLock()
	defer m.RUnlock()

	periods := make([]interface{}, 0)

	for ticker, tickPeriod := range m.hash {
		periods = append(periods, tickPeriod)
		m.publish(tickPeriod)
		// reset hash
		m.hash[ticker] = &Period{Ticker: ticker}
	}

	m.persist(periods)
}

func (m *PeriodHash) publish(tickPeriod *Period) {
	// fmt.Println("TICKER_SERVICE: ", tickPeriod)
}

func (m *PeriodHash) persist(l []interface{}) {
	// put it all into mongodb
	c := m.db.C("ticks")
	err := c.Insert(l...)
	if err != nil {
		panic(err)
	}
}
