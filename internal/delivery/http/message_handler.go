package http

import (
	"encontro/internal/delivery/http/dto"
	"encontro/internal/domain/entity"
	"encontro/internal/usecase"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MessageHandler обрабатывает HTTP-запросы для сообщений
type MessageHandler struct {
	messageUseCase MessageUseCaseInterface
}

// NewMessageHandler создает новый экземпляр MessageHandler
func NewMessageHandler(messageUseCase MessageUseCaseInterface) *MessageHandler {
	return &MessageHandler{
		messageUseCase: messageUseCase,
	}
}

// CreateMessage обрабатывает создание нового сообщения
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var req dto.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	// Получаем senderID из заголовка или из тела запроса
	senderID := c.GetHeader("X-User-ID")
	if senderID == "" {
		if req.SenderID != "" {
			senderID = req.SenderID
		} else if req.UserID != "" {
			senderID = req.UserID
		}
	}
	if senderID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	msg, err := h.messageUseCase.CreateMessage(c.Request.Context(), req.RoomID, req.Content, senderID)
	if err != nil {
		switch err {
		case usecase.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case usecase.ErrRoomNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			fmt.Printf("Error creating message: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("internal server error: %v", err)})
		}
		return
	}

	c.JSON(http.StatusCreated, mapMessageToResponse(msg))
}

// GetMessages обрабатывает получение списка сообщений с пагинацией
func (h *MessageHandler) GetMessages(c *gin.Context) {
	var params entity.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination parameters: " + err.Error()})
		return
	}

	// Устанавливаем значения по умолчанию
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 20
	}

	roomID := c.Query("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room ID is required"})
		return
	}

	response, err := h.messageUseCase.GetMessagesByRoomID(c.Request.Context(), roomID, params)
	if err != nil {
		switch err {
		case usecase.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case usecase.ErrRoomNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// isValidUUID проверяет, что строка является валидным UUID
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// GetMessage обрабатывает получение сообщения по ID
func (h *MessageHandler) GetMessage(c *gin.Context) {
	id := c.Param("messageId")
	if id == "" || !isValidUUID(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}

	msg, err := h.messageUseCase.GetMessageByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case usecase.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case usecase.ErrMessageNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, mapMessageToResponse(msg))
}

// UpdateMessage обрабатывает обновление сообщения
func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	id := c.Param("messageId")
	if id == "" || !isValidUUID(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}

	var req dto.UpdateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	// Получаем существующее сообщение
	existingMsg, err := h.messageUseCase.GetMessageByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case usecase.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case usecase.ErrMessageNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// Обновляем только контент
	existingMsg.Content = req.Content
	existingMsg.UpdatedAt = time.Now()

	// Сохраняем обновленное сообщение
	if err := h.messageUseCase.UpdateMessage(c.Request.Context(), existingMsg); err != nil {
		switch err {
		case usecase.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case usecase.ErrMessageNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// Получаем обновленное сообщение для ответа
	updatedMsg, err := h.messageUseCase.GetMessageByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get updated message"})
		return
	}

	c.JSON(http.StatusOK, mapMessageToResponse(updatedMsg))
}

// DeleteMessage обрабатывает удаление сообщения
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	id := c.Param("messageId")
	if id == "" || !isValidUUID(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}

	// Сначала проверяем существование сообщения
	if _, err := h.messageUseCase.GetMessageByID(c.Request.Context(), id); err != nil {
		switch err {
		case usecase.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case usecase.ErrMessageNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// Затем удаляем его
	if err := h.messageUseCase.DeleteMessage(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete message"})
		return
	}

	c.Status(http.StatusNoContent)
}

// mapMessageToResponse преобразует сущность Message в DTO MessageResponse
func mapMessageToResponse(msg *entity.Message) dto.MessageResponse {
	return dto.MessageResponse{
		ID:        msg.ID,
		Content:   msg.Content,
		RoomID:    msg.RoomID,
		SenderID:  msg.SenderID,
		CreatedAt: msg.CreatedAt.Format(time.RFC3339),
		UpdatedAt: msg.UpdatedAt.Format(time.RFC3339),
	}
}
