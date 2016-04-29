package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"text/template"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var index = template.Must(template.ParseFiles("./web_client/index.html"))
var css = template.Must(template.ParseFiles("./web_client/chart.css"))
var js = template.Must(template.ParseFiles("./web_client/chart.js"))
var dataJs = template.Must(template.ParseFiles("./web_client/data.js"))

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/chart.css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
			css.Execute(w, nil)

		case "/chart.js":
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			js.Execute(w, nil)

		case "/data.js":
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			dataJs.Execute(w, nil)

		case "/":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			index.Execute(w, nil)
		}

	})

	http.HandleFunc("/msg", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status 200"))
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}
		fmt.Println("websocket connected")
		fmt.Println(ws)
	})

	http.ListenAndServe(":8004", nil)
}
