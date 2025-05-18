package entity

import "time"

// Message представляет сообщение, отправленное в чате (например, в комнате).
type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	RoomID    string    `json:"room_id"`
	SenderID  string    `json:"sender_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
