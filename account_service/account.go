package main

import (
	"net/http"
	"fmt"
 	// "time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		// time.Sleep(5000 * time.Millisecond) 
		// Ive discovered new problems...
		// if this request blocks in response,
		// the market_service cant keep handling trades
		// this is where node.js shines (non blocking)

		// decoder := json.NewDecoder(req.Body)
		fmt.Println("hello account")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})
	http.ListenAndServe(":8000", nil)
}