package http

import (
	"encontro/internal/delivery/http/dto"
	"encontro/internal/domain/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoomHandler обрабатывает HTTP-запросы для комнат
type RoomHandler struct {
	roomUseCase RoomUseCaseInterface
}

// NewRoomHandler создает новый экземпляр RoomHandler
func NewRoomHandler(roomUseCase RoomUseCaseInterface) *RoomHandler {
	return &RoomHandler{roomUseCase: roomUseCase}
}

// CreateRoom обрабатывает создание новой комнаты
// @Summary Создать комнату
// @Accept json
// @Produce json
// @Param request body dto.CreateRoomRequest true "Параметры комнаты"
// @Success 201 {object} entity.Room
// @Failure 400 {object} gin.H
// @Router /rooms [post]
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req dto.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	room, err := h.roomUseCase.CreateRoom(c.Request.Context(), req.Name, req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания комнаты: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, room)
}

// GetRoom обрабатывает получение комнаты по ID
// @Summary Получить комнату по ID
// @Produce json
// @Param id path string true "ID комнаты"
// @Success 200 {object} entity.Room
// @Failure 404 {object} gin.H
// @Router /rooms/{id} [get]
func (h *RoomHandler) GetRoom(c *gin.Context) {
	id := c.Param("id")
	room, err := h.roomUseCase.GetRoom(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	c.JSON(http.StatusOK, room)
}

// ListRooms обрабатывает получение списка комнат с пагинацией
// @Summary Получить список всех комнат
// @Produce json
// @Success 200 {array} entity.Room
// @Router /rooms [get]
func (h *RoomHandler) ListRooms(c *gin.Context) {
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

	response, err := h.roomUseCase.GetRooms(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteRoom обрабатывает удаление комнаты
// @Summary Удалить комнату
// @Param id path string true "ID комнаты"
// @Success 204
// @Failure 404 {object} gin.H
// @Router /rooms/{id} [delete]
func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	if err := h.roomUseCase.DeleteRoom(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
