package websocket

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu    sync.RWMutex
	conns map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		conns: make(map[string]*websocket.Conn),
	}
}

func (h *Hub) Register(clientID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.conns[clientID] = conn
}

func (h *Hub) Unregister(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.conns, clientID)
}

func (h *Hub) SendTo(clientID string, messageType int, data []byte) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conn, ok := h.conns[clientID]

	if !ok {
		return fmt.Errorf("client %s not found", clientID)
	}

	return conn.WriteMessage(messageType, data)
}

func (h *Hub) SendJSONTo(clientID string, payload interface{}) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conn, ok := h.conns[clientID]

	if !ok {
		return fmt.Errorf("client %s not found", clientID)
	}

	return conn.WriteJSON(payload)
}
