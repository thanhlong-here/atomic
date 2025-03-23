package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceID := r.URL.Query().Get("service_id")
		if serviceID == "" {
			http.Error(w, "Missing service_id", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}

		client := &Client{
			ID:     serviceID,
			Conn:   conn,
			Send:   make(chan Message, 256),
			Topics: make(map[string]bool),
		}

		hub.Register <- client
		go client.ReadPump(hub)
	}
}
