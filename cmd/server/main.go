package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	httphandler "encontro/internal/delivery/http"
	"encontro/internal/delivery/websocket"
	"encontro/internal/domain/service"
	inmemRepo "encontro/internal/infrastructure/repository"
	"encontro/internal/usecase"
)

func loadTLSConfig() (*tls.Config, error) {
	certPath := filepath.Join("certs", "localhost+2.pem")
	keyPath := filepath.Join("certs", "localhost+2-key.pem")

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certs: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}

func main() {
	// Создаем конфигурацию для репозиториев
	repoConfig := inmemRepo.NewInMemoryConfig(false)

	// Инициализация репозитория
	roomRepo := inmemRepo.NewInMemoryRoomRepository(repoConfig)

	// Инициализация UUID генератора
	uuidGen := service.NewGoogleUUIDGenerator()

	// Инициализация usecase для сообщений (messageRepo должен быть создан ранее)
	messageRepo := inmemRepo.NewInMemoryMessageRepository(repoConfig)
	messageUseCase := usecase.NewMessageUseCase(messageRepo, uuidGen)

	// Инициализация usecase
	roomUseCase := usecase.NewRoomUseCase(roomRepo, uuidGen)

	// Инициализация обработчиков
	wsHandler := websocket.NewHandler(roomUseCase)
	roomHandler := httphandler.NewRoomHandler(roomUseCase)

	// Инициализация обработчика сообщений
	messageHandler := httphandler.NewMessageHandler(messageUseCase)

	// Настройка маршрутизатора
	router := gin.Default()

	// Статические файлы
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	// API группа
	api := router.Group("/api")
	{
		// WebSocket endpoint
		api.GET("/ws/:room", wsHandler.HandleWebSocket)

		// CRUD endpoints для комнат
		rooms := api.Group("/rooms")
		{
			rooms.POST("", roomHandler.CreateRoom)       // POST /api/rooms - создать комнату
			rooms.GET("", roomHandler.ListRooms)         // GET /api/rooms - список всех комнат
			rooms.GET("/:id", roomHandler.GetRoom)       // GET /api/rooms/:id - получить комнату по ID
			rooms.DELETE("/:id", roomHandler.DeleteRoom) // DELETE /api/rooms/:id - удалить комнату
		}

		// API для сообщений
		messages := api.Group("/messages")
		{
			messages.POST("", messageHandler.CreateMessage)
			messages.GET("", messageHandler.GetMessages)
			messages.GET("/:id", messageHandler.GetMessage)
			messages.DELETE("/:id", messageHandler.DeleteMessage)
		}
	}

	// Загрузка TLS конфигурации
	tlsConfig, err := loadTLSConfig()
	if err != nil {
		log.Fatalf("Failed to load TLS config: %v", err)
	}

	// Создание HTTP сервера
	server := &http.Server{
		Addr:      ":8443",
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	// Канал для получения сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Printf("Starting server at https://localhost:8443")
		log.Printf("Please note:")
		log.Printf("  * Note the HTTPS in the URL; there is no HTTP -> HTTPS redirect.")
		log.Printf("  * You'll need to accept the invalid TLS certificate as it is self-signed.")
		log.Printf("  * Some browsers or OSs may not allow the webcam to be used by multiple pages at once.")

		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-stop
	log.Println("Shutting down server...")

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
