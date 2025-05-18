package http

import (
	"encontro/internal/domain/entity"
	"encontro/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MessageHandler обрабатывает HTTP-запросы для сообщений
type MessageHandler struct {
	messageUseCase *usecase.MessageUseCase
}

// NewMessageHandler создает новый экземпляр MessageHandler
func NewMessageHandler(messageUseCase *usecase.MessageUseCase) *MessageHandler {
	return &MessageHandler{
		messageUseCase: messageUseCase,
	}
}

// CreateMessage обрабатывает создание нового сообщения
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Param("id")
	senderID := c.GetString("user_id") // Предполагается, что middleware аутентификации установил user_id

	msg, err := h.messageUseCase.CreateMessage(c.Request.Context(), roomID, req.Content, senderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}

// GetMessages обрабатывает получение списка сообщений с пагинацией
func (h *MessageHandler) GetMessages(c *gin.Context) {
	var params entity.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем значения по умолчанию
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 20
	}

	roomID := c.Param("id")
	response, err := h.messageUseCase.GetMessagesByRoomID(c.Request.Context(), roomID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMessage обрабатывает получение сообщения по ID
func (h *MessageHandler) GetMessage(c *gin.Context) {
	id := c.Param("messageId")
	msg, err := h.messageUseCase.GetMessageByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if msg == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}

	c.JSON(http.StatusOK, msg)
}

// DeleteMessage обрабатывает удаление сообщения
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	id := c.Param("messageId")
	if err := h.messageUseCase.DeleteMessage(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
