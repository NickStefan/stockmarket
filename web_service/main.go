package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/msg", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})

	http.ListenAndServe(":8004", nil)
}
