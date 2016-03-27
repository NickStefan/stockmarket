package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

type Trade struct {
	Actor string `json:"actor"`
	Shares int `json:"shares"`
	Price float64 `json:"price"`
	Intent string `json:"intent"`
	Kind string `json:"kind"`
	State  string `json:"state"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		var t Trade
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
		fmt.Println(t)
	})
	http.ListenAndServe(":8000", nil)
}