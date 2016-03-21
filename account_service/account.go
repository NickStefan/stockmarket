package main

import (
	"net/http"
	"fmt"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){

		// if this request blocks in response,
		// the market_service cant keep handling trades

		// that should actually be okay. we probably dont want to do any trading
		// if the order fufillment is down!

		// but, we dont want to be limited by each individual trade being an HTTP
		// round trip. what if we used batching?

		// what if we batch all of the trades in the market_service every 100ms
		// we then every 100ms send the trades to the account_service?

		// decoder := json.NewDecoder(req.Body)
		fmt.Println("hello account")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})
	http.ListenAndServe(":8000", nil)
}