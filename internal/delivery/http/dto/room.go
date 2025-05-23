package dto

// CreateRoomRequest представляет запрос на создание комнаты
type CreateRoomRequest struct {
	Name string `json:"name" validate:"required"`
}

// UpdateRoomRequest представляет запрос на обновление комнаты
type UpdateRoomRequest struct {
	Name string `json:"name" validate:"required"`
}

// RoomResponse представляет ответ с информацией о комнате
type RoomResponse struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Clients   []ClientInfo `json:"clients"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}

// ClientInfo представляет информацию о клиенте
type ClientInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// RoomListResponse представляет ответ со списком комнат
type RoomListResponse struct {
	Rooms []RoomResponse `json:"rooms"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}
