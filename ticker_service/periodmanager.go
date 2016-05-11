package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"sync"
	"time"
)

type Period struct {
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Open   float64   `json:"open"`
	Close  float64   `json:"close"`
	Volume int       `json:"volume"`
	Ticker string    `json:"ticker"`
	Date   time.Time `json:"date"`
}

type PeriodManager struct {
	tickers   []string
	hash      *PeriodHash
	db        *mgo.Database
	publisher func(*Period)
	sync.RWMutex
}

func NewPeriodManager(tickers []string, hash *PeriodHash) *PeriodManager {
	for _, ticker := range tickers {
		hash.set(ticker, &Period{Ticker: ticker, Date: time.Now()})
	}

	return &PeriodManager{
		hash:    hash,
		tickers: tickers,
	}
}

func (m *PeriodManager) setDB(db *mgo.Database) {
	m.db = db
}

func (m *PeriodManager) setPublisher(p func(*Period)) {
	m.publisher = p
}

func (m *PeriodManager) add(t Trade) {
	m.Lock()
	defer m.Unlock()

	period := m.hash.get(t.Ticker)

	if period.Volume == 0 {
		period.Open = t.Price
		period.High = t.Price
		period.Low = t.Price
		period.Close = t.Price
	}

	if period.Low > t.Price {
		period.Low = t.Price
	}

	if period.High < t.Price {
		period.High = t.Price
	}

	period.Volume = period.Volume + t.Shares
	period.Close = t.Price

	m.hash.set(t.Ticker, period)
}

func (m *PeriodManager) Persist() {
	m.RLock()
	defer m.RUnlock()

	periods := make([]interface{}, 0)

	for _, ticker := range m.tickers {
		periods = append(periods, m.hash.get(ticker))
		// reset hash
		m.hash.set(ticker, &Period{Ticker: ticker, Date: time.Now()})
	}

	m.persist(periods)
}

func (m *PeriodManager) persist(l []interface{}) {
	// put it all into mongodb
	c := m.db.C("ticks")
	err := c.Insert(l...)
	if err != nil {
		fmt.Println("TODO: fault tolerance needed; ", err)
	}
}

func (m *PeriodManager) Publish() {
	m.RLock()
	defer m.RUnlock()

	for _, ticker := range m.tickers {
		m.publish(m.hash.get(ticker))
		// reset hash
		m.hash.set(ticker, &Period{Ticker: ticker, Date: time.Now()})
	}
}

func (m *PeriodManager) publish(tickPeriod *Period) {
	m.publisher(tickPeriod)
}
