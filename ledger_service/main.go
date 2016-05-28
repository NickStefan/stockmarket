package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var mutex sync.Mutex
	dataStore := make(map[string]*Ledger)

	http.HandleFunc("/fill", func(w http.ResponseWriter, r *http.Request) {
		var payload [2]Trade
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&payload)
		if err != nil {
			fmt.Println("TODO: ledger_service fault tolerance needed; ", err)
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

	http.ListenAndServe(":8080", nil)
}
