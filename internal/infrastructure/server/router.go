package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Config содержит конфигурацию для сервера
type Config struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}
}

// NewRouter создает новый маршрутизатор с настройками CORS
func NewRouter(config *Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	if config == nil {
		config = DefaultConfig()
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.AllowOrigins
	corsConfig.AllowMethods = config.AllowMethods
	corsConfig.AllowHeaders = config.AllowHeaders
	corsConfig.AllowCredentials = config.AllowCredentials

	router.Use(cors.New(corsConfig))
	return router
}
