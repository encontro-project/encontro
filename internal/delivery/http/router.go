package http

import (
	"encontro/internal/delivery/http/middleware"
	"encontro/internal/usecase"

	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает маршруты приложения
func SetupRouter(roomUseCase *usecase.RoomUseCase, messageUseCase *usecase.MessageUseCase) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// Хендлеры
	roomHandler := NewRoomHandler(roomUseCase)
	messageHandler := NewMessageHandler(messageUseCase)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Маршруты для комнат
		rooms := v1.Group("/rooms")
		{
			rooms.POST("", roomHandler.CreateRoom)
			rooms.GET("", roomHandler.ListRooms)
			rooms.GET("/:id", roomHandler.GetRoom)
			rooms.DELETE("/:id", roomHandler.DeleteRoom)
		}

		// Маршруты для сообщений
		messages := v1.Group("/rooms/:roomId/messages")
		{
			messages.POST("", messageHandler.CreateMessage)
			messages.GET("", messageHandler.GetMessages) // Новый маршрут для пагинированного списка
			messages.GET("/:id", messageHandler.GetMessage)
			messages.DELETE("/:id", messageHandler.DeleteMessage)
		}
	}

	return router
}
