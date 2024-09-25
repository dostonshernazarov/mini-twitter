package websocket

import (
	"log"
	"net/http"

	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/kafka"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	kafka.AddClient(conn)
	defer kafka.RemoveClient(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket closed:", err)
			return
		}
	}
}
