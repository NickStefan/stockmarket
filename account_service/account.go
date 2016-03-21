package main

import (
	"net/http"
	"fmt"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		// decoder := json.NewDecoder(req.Body)
		fmt.Println("hello account")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})
	http.ListenAndServe(":8000", nil)
}