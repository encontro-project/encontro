package dto

// CreateMessageRequest представляет запрос на создание сообщения
type CreateMessageRequest struct {
	Content  string `json:"content" binding:"required"`
	RoomID   string `json:"room_id" binding:"required"`
	SenderID string `json:"sender_id"`
	UserID   string `json:"user_id"`
}

// UpdateMessageRequest представляет запрос на обновление сообщения
type UpdateMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// MessageResponse представляет ответ с информацией о сообщении
type MessageResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	RoomID    string `json:"room_id"`
	SenderID  string `json:"sender_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// MessageListResponse представляет ответ со списком сообщений
type MessageListResponse struct {
	Messages []MessageResponse `json:"messages"`
	Page     int               `json:"page"`
	Size     int               `json:"size"`
	Total    int64             `json:"total"`
}
