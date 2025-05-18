package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"encontro/internal/domain/entity"
	"encontro/internal/usecase"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // В продакшене нужно настроить правильную проверку origin
	},
}

// Message представляет собой сообщение WebSocket
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Handler обрабатывает WebSocket соединения
type Handler struct {
	roomUseCase *usecase.RoomUseCase
	mu          sync.RWMutex
}

// NewHandler создает новый экземпляр Handler
func NewHandler(roomUseCase *usecase.RoomUseCase) *Handler {
	return &Handler{
		roomUseCase: roomUseCase,
	}
}

// HandleWebSocket обрабатывает входящие WebSocket соединения
func (h *Handler) HandleWebSocket(c *gin.Context) {
	roomID := c.Param("room")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room ID is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Создаем клиента
	client := &entity.Client{
		ID:       fmt.Sprintf("client_%d", time.Now().UnixNano()),
		RoomID:   roomID,
		Username: fmt.Sprintf("User_%d", time.Now().UnixNano()),
	}

	// Добавляем клиента в комнату
	if err := h.roomUseCase.AddClientToRoom(c.Request.Context(), roomID, client); err != nil {
		log.Printf("Failed to add client to room: %v", err)
		return
	}
	defer h.roomUseCase.RemoveClientFromRoom(c.Request.Context(), roomID, client.ID)

	// Канал для получения сообщений от клиента
	messageChan := make(chan Message)
	// Канал для сигнала завершения
	done := make(chan struct{})

	// Горутина для чтения сообщений
	go func() {
		defer close(messageChan)
		for {
			var msg Message
			if err := conn.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket read error: %v", err)
				}
				return
			}
			messageChan <- msg
		}
	}()

	// Обработка сообщений
	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				return
			}
			if err := h.handleMessage(c.Request.Context(), conn, client, msg); err != nil {
				log.Printf("Failed to handle message: %v", err)
				return
			}
		case <-done:
			return
		}
	}
}

// handleMessage обрабатывает входящие сообщения
func (h *Handler) handleMessage(ctx context.Context, conn *websocket.Conn, client *entity.Client, msg Message) error {
	switch msg.Type {
	case "join":
		// Обработка присоединения к комнате
		room, err := h.roomUseCase.GetRoom(ctx, client.RoomID)
		if err != nil {
			return fmt.Errorf("failed to get room: %w", err)
		}

		// Отправляем информацию о других клиентах в комнате
		clients := room.GetClients()
		response := struct {
			Type    string           `json:"type"`
			Payload []*entity.Client `json:"payload"`
		}{
			Type:    "room_state",
			Payload: clients,
		}

		if err := conn.WriteJSON(response); err != nil {
			return fmt.Errorf("failed to send room state: %w", err)
		}

	case "offer", "answer", "ice-candidate":
		// Пересылаем сигнальные сообщения другим клиентам
		room, err := h.roomUseCase.GetRoom(ctx, client.RoomID)
		if err != nil {
			return fmt.Errorf("failed to get room: %w", err)
		}

		// Отправляем сообщение всем клиентам в комнате, кроме отправителя
		for _, c := range room.GetClients() {
			if c.ID != client.ID {
				// TODO(#2): Реализовать отправлениe сигнального сообщения конкретному клиенту
				//   Необходимо организовать хранение активных WebSocket-соединений,
				//   чтобы можно было находить нужное соединение по ID клиента и пересылать ему сообщение
			}
		}

	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}

	return nil
}
