package http

import (
	"encontro/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает маршруты приложения
func SetupRouter(
	roomUseCase RoomUseCaseInterface,
	messageUseCase MessageUseCaseInterface,
	userUseCase UserSummaryUseCaseInterface,
	wsHandler WebSocketHandlerInterface,
) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// Обработка OPTIONS-запросов для всех путей
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	// Хендлеры
	roomHandler := NewRoomHandler(roomUseCase)
	messageHandler := NewMessageHandler(messageUseCase)
	userHandler := NewUserSummaryHandler(userUseCase)

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
		messages := v1.Group("/messages")
		{
			messages.POST("", messageHandler.CreateMessage)
			messages.GET("", messageHandler.GetMessages)
			messages.GET("/:messageId", messageHandler.GetMessage)
			messages.PUT("/:messageId", messageHandler.UpdateMessage)
			messages.DELETE("/:messageId", messageHandler.DeleteMessage)
		}

		// Маршруты для пользователей
		users := v1.Group("/user")
		{
			users.Use(middleware.UserIDMiddleware())
			users.GET("/:id", userHandler.GetUserSummary)
		}
	}

	// WebSocket маршрут (без версионирования)
	router.GET("/api/ws/:room", wsHandler.HandleWebSocket)

	// Healthcheck endpoint
	router.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})

	return router
}
