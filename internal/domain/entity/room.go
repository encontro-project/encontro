package entity

import (
	"sync"
	"time"
)

// Room представляет собой комнату для видеочата
type Room struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Clients   []*Client `json:"clients"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	mu        sync.RWMutex
}

// Client представляет собой клиента в комнате
type Client struct {
	ID       string `json:"id"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
	mu       sync.RWMutex
}

// NewRoom создает новую комнату
func NewRoom(id, name string) *Room {
	now := time.Now()
	return &Room{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
		Clients:   make([]*Client, 0),
	}
}

// AddClient добавляет клиента в комнату
func (r *Room) AddClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Clients = append(r.Clients, client)
	r.UpdatedAt = time.Now()
}

// RemoveClient удаляет клиента из комнаты
func (r *Room) RemoveClient(clientID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, client := range r.Clients {
		if client.ID == clientID {
			r.Clients = append(r.Clients[:i], r.Clients[i+1:]...)
			break
		}
	}
	r.UpdatedAt = time.Now()
}

// GetClients возвращает список клиентов в комнате
func (r *Room) GetClients() []*Client {
	return r.Clients
}

// GetClient возвращает клиента по ID
func (r *Room) GetClient(clientID string) (*Client, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, client := range r.Clients {
		if client.ID == clientID {
			return client, true
		}
	}
	return nil, false
}
