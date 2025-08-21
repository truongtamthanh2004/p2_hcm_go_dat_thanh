package config

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Client struct {
	UserID uint
	Conn   *websocket.Conn
}

type Hub struct {
	clients map[uint]*Client
	lock    sync.RWMutex
}

var hub = Hub{
	clients: make(map[uint]*Client),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWS(w http.ResponseWriter, r *http.Request, userID uint) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{UserID: userID, Conn: conn}

	hub.lock.Lock()
	hub.clients[userID] = client
	hub.lock.Unlock()

	log.Printf("User %d connected\n", userID)

	// Listen
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			hub.lock.Lock()
			delete(hub.clients, userID)
			hub.lock.Unlock()
			conn.Close()
			log.Printf("User %d disconnected\n", userID)
			break
		}
	}
}

func SendToUser(userID uint, message interface{}) {
	hub.lock.RLock()
	client, ok := hub.clients[userID]
	hub.lock.RUnlock()
	if ok {
		data, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message for user %d: %v\n", userID, err)
			return
		}
		if err := client.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Error sending message to user %d: %v\n", userID, err)
		}
	}
}
