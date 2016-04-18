package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Trade struct {
	Actor  string  `json:"actor"`
	Shares int     `json:"shares"`
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
	Intent string  `json:"intent"`
	Kind   string  `json:"kind"`
	State  string  `json:"state"`
}

type Asset struct {
	ticker string
	shares int
}

type Ledger struct {
	name   string
	cash   float64
	assets map[string]*Asset
}

func (a *Ledger) buy(t Trade) {
	if a.assets[t.Ticker] == nil {
		a.assets[t.Ticker] = &Asset{}
	}
	a.cash = a.cash - (float64(t.Shares) * t.Price)
	asset := a.assets[t.Ticker]
	asset.ticker = t.Ticker
	asset.shares = asset.shares + t.Shares
}

func (a *Ledger) sell(t Trade) {
	if a.assets[t.Ticker] == nil {
		a.assets[t.Ticker] = &Asset{}
	}
	a.cash = a.cash + (float64(t.Shares) * t.Price)
	asset := a.assets[t.Ticker]
	asset.ticker = t.Ticker
	asset.shares = asset.shares - t.Shares
}

func processTrade(data map[string]*Ledger, t Trade, o Trade) {
	if data[t.Actor] == nil {
		data[t.Actor] = &Ledger{name: t.Actor, cash: 0, assets: make(map[string]*Asset)}
	}
	if data[o.Actor] == nil {
		data[o.Actor] = &Ledger{name: o.Actor, cash: 0, assets: make(map[string]*Asset)}
	}

	if t.Intent == "BUY" {
		data[t.Actor].buy(t)

	} else if t.Intent == "SELL" {
		data[t.Actor].sell(t)
	}
	
	if o.Intent == "BUY" {
		data[o.Actor].buy(o)

	} else if o.Intent == "SELL" {
		data[o.Actor].sell(o)
	}

}

func main() {
	var mutex sync.Mutex
	dataStore := make(map[string]*Ledger)

	http.HandleFunc("/fill", func(w http.ResponseWriter, r *http.Request) {
		var payload [2]Trade
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&payload)
		if err != nil {
			fmt.Println("ERR: LEDGER_SERVICE")
			panic(err)
		}

		mutex.Lock()
		defer mutex.Unlock()
		processTrade(dataStore, payload[0], payload[1])

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})

	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		for name, ledger := range dataStore {
			fmt.Println("LEDGER_SERVICE: ", name, ledger)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})

	http.ListenAndServe(":8002", nil)
}
