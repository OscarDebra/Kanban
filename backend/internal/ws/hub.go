package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	send    chan []byte
	boardID int
	userID  int
}

type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*Client]bool)}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	h.clients[c] = true
	h.mu.Unlock()
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	delete(h.clients, c)
	close(c.send)
	h.mu.Unlock()
}

func (h *Hub) BroadcastToBoard(boardID int, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.boardID == boardID {
			select {
			case client.send <- message:
			default:
				go h.Unregister(client)
			}
		}
	}
}