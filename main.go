package main

import (
	"atomic/internal/ws"
	"net/http"
)

func main() {
	hub := ws.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", ws.HandleWS(hub))
	http.ListenAndServe(":8082", nil)
}
