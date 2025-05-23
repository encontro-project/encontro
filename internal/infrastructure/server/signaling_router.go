package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SignalingConfig содержит конфигурацию для signaling сервера
type SignalingConfig struct {
	AllowOrigins     []string
	AllowCredentials bool
}

// DefaultSignalingConfig возвращает конфигурацию по умолчанию для signaling сервера
func DefaultSignalingConfig() *SignalingConfig {
	return &SignalingConfig{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175"},
		AllowCredentials: true,
	}
}

// NewSignalingRouter создает новый маршрутизатор для signaling сервера
// с минимальными настройками CORS, необходимыми для WebSocket
func NewSignalingRouter(config *SignalingConfig) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	if config == nil {
		config = DefaultSignalingConfig()
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.AllowOrigins
	corsConfig.AllowMethods = []string{"GET", "OPTIONS"} // Только GET для WebSocket
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "Cache-Control", "X-Requested-With"}
	corsConfig.AllowCredentials = config.AllowCredentials
	corsConfig.MaxAge = 86400 // 24 часа

	router.Use(cors.New(corsConfig))

	// Добавляем обработчик OPTIONS для всех путей
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	return router
}
