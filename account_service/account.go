package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

type Trade struct {
	Actor string `json:"actor"`
	Shares int 	`json:"shares"`
	Ticker string `json:"ticker"`
	Price float64 `json:"price"`
	Intent string `json:"intent"`
	Kind string `json:"kind"`
	State  string `json:"state"`
}

type Account struct {
	name string
	cash float64
	assets map[string]*Asset
}

type Asset struct {
	ticker string
	shares int
}

func processTrade(data map[string]*Account, t Trade){
	if data[t.Actor] == nil {
		data[t.Actor] = &Account{
			name: t.Actor, cash: 10000,
		}
	}
}


func main() {

	dataStore := make(map[string]*Account)

	http.HandleFunc("/fill", func(w http.ResponseWriter, r *http.Request){
		var t Trade
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))

		processTrade(dataStore, t)
		fmt.Println(t)
	})

	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request){
		for name, account := range dataStore {
			fmt.Println(name, account)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})

	http.ListenAndServe(":8000", nil)
}