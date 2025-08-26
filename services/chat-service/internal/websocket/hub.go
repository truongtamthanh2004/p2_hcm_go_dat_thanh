package websocket

import (
	"chat-service/internal/constant"
	"chat-service/internal/dto"
	"chat-service/internal/model"
	"chat-service/internal/usecase"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	clients    map[uint]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *model.ChatMessage
	usecase    usecase.ChatUsecase
}

func NewHub(chatUsecase usecase.ChatUsecase) *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *model.ChatMessage),
		usecase:    chatUsecase,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.UserID] = client

		case client := <-h.unregister:
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}

		case msg := <-h.broadcast:
			// Save message to DB
			if err := h.usecase.SaveMessage(context.Background(), msg); err != nil {
				log.Printf("failed to save message: %v", err)
				continue
			}

			// Send to receiver if online
			if receiver, ok := h.clients[msg.ReceiverID]; ok {
				dto := dto.ChatMessageDTO{
					ID:         msg.ID,
					SenderID:   msg.SenderID,
					ReceiverID: msg.ReceiverID,
					Content:    msg.Content,
					CreatedAt:  msg.CreatedAt,
				}

				data, err := json.Marshal(dto)
				if err != nil {
					log.Printf("failed to marshal message dto: %v", err)
					continue
				}
				receiver.Send <- data
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) error {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		return errors.New(constant.ErrUserIDRequired)
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return errors.New(constant.ErrInvalidUserID)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.New(constant.ErrUpgradeFailed)
	}

	client := &Client{
		UserID: uint(userID),
		Conn:   conn,
		Send:   make(chan []byte),
	}

	h.register <- client

	go client.writePump()
	go client.readPump(h)
	return nil
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msgData, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg model.ChatMessage
		if err := json.Unmarshal(msgData, &msg); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}

		if len(msg.Content) == 0 || len(msg.Content) > 1000 {
			log.Printf("Invalid message content length from user %d", c.UserID)
			continue
		}

		if msg.ReceiverID == 0 {
			log.Printf("Invalid receiver ID from user %d", c.UserID)
			continue
		}

		msg.SenderID = c.UserID
		h.broadcast <- &msg
	}
}

func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()

	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("write error to user %d: %v", c.UserID, err)
			break
		}
	}
}
