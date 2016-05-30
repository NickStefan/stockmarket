package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"gopkg.in/redsync.v1"
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
	env       string
	period    string
	tickers   []string
	hash      *PeriodHash
	db        *mgo.Database
	publisher func(*Period)
	redsync   *redsync.Redsync
	lockMap   map[string]*Locker
}

func NewPeriodManager(pool *redis.Pool, hash *PeriodHash, period string) *PeriodManager {
	return &PeriodManager{
		period:  period,
		lockMap: make(map[string]*Locker),
		redsync: redsync.New([]redsync.Pool{pool}),
		hash:    hash,
	}
}

func (m *PeriodManager) setDB(db *mgo.Database) {
	m.db = db
}

func (m *PeriodManager) setPublisher(p func(*Period)) {
	m.publisher = p
}

func (m *PeriodManager) setEnv(env string) {
	m.env = env
}

func (m *PeriodManager) setTickers(tickers []string) {
	m.tickers = tickers
}

func (m *PeriodManager) initPeriods() error {
	for _, ticker := range m.tickers {
		locker := m.getLocker(ticker)
		err := locker.Lock()
		if err != nil {
			return err
		}
		if nil == m.hash.get(ticker) {
			m.hash.set(ticker, &Period{Ticker: ticker, Date: time.Now()})
		}
		locker.Unlock()
	}
	return nil
}

func (m *PeriodManager) getLocker(ticker string) *Locker {
	if nil != m.lockMap[ticker] {
		return m.lockMap[ticker]
	} else {
		redLockMutex := m.redsync.NewMutex("ticker_service" + m.period + ticker)
		redsync.SetRetryDelay(5 * time.Millisecond).Apply(redLockMutex)
		redsync.SetExpiry(500 * time.Millisecond).Apply(redLockMutex)
		redsync.SetTries(50).Apply(redLockMutex)

		m.lockMap[ticker] = &Locker{
			name:    "ticker_service" + m.period + ticker,
			env:     m.env,
			mutLock: &sync.Mutex{},
			redLock: redLockMutex,
		}
		return m.lockMap[ticker]
	}
}

func (m *PeriodManager) add(t Trade) error {
	locker := m.getLocker(t.Ticker)
	err := locker.Lock()
	if err != nil {
		return err
	}
	defer locker.Unlock()

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
	return nil
}

func (m *PeriodManager) Persist() error {

	periods := make([]interface{}, 0)

	for _, ticker := range m.tickers {
		locker := m.getLocker(ticker)
		err := locker.Lock()
		if err != nil {
			return err
		}
		periods = append(periods, m.hash.get(ticker))
		// reset hash
		m.hash.set(ticker, &Period{Ticker: ticker, Date: time.Now()})
		locker.Unlock()
	}

	m.persist(periods)
	return nil
}

func (m *PeriodManager) persist(l []interface{}) {
	// put it all into mongodb
	c := m.db.C("ticks")
	err := c.Insert(l...)
	if err != nil {
		fmt.Println("ticker_service: periodmanager mongodb", err, l)
	}
}

func (m *PeriodManager) Publish() error {
	for _, ticker := range m.tickers {
		locker := m.getLocker(ticker)
		err := locker.Lock()
		if err != nil {
			return err
		}
		m.publish(m.hash.get(ticker))
		// reset hash
		m.hash.set(ticker, &Period{Ticker: ticker, Date: time.Now()})
		locker.Unlock()
	}
	return nil
}

func (m *PeriodManager) publish(tickPeriod *Period) {
	go m.publisher(tickPeriod)
}
