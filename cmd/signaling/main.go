package main

import (
	"log"
	"os"
	"path/filepath"

	"encontro/internal/delivery/websocket"
	"encontro/internal/domain/service"
	inmemRepo "encontro/internal/infrastructure/repository"
	"encontro/internal/infrastructure/server"
	"encontro/internal/usecase"
)

func main() {
	// Получаем конфигурацию из переменных окружения
	port := os.Getenv("SIGNALING_PORT")
	if port == "" {
		port = "8443" // порт по умолчанию для signaling
	}
	certPath := os.Getenv("CERT_PATH")
	if certPath == "" {
		certPath = "certs" // путь по умолчанию
	}

	// Инициализация репозитория и use case для комнат
	repoConfig := inmemRepo.NewInMemoryConfig(false)
	roomRepo, err := inmemRepo.NewInMemoryRoomRepository(repoConfig)
	if err != nil {
		log.Fatalf("Failed to create room repository: %v", err)
	}
	uuidGen := service.NewGoogleUUIDGenerator()
	roomUseCase := usecase.NewRoomUseCase(roomRepo, uuidGen)

	// Инициализация WebSocket обработчика
	wsHandler := websocket.NewHandler(roomUseCase)

	// Настройка конфигурации signaling сервера
	signalingConfig := server.DefaultSignalingConfig()

	// Настройка маршрутизатора
	router := server.NewSignalingRouter(signalingConfig)

	// WebSocket endpoint
	router.GET("/api/ws/:room", wsHandler.HandleWebSocket)

	// Пути к сертификатам
	certFile := filepath.Join(certPath, "localhost+2.pem")
	keyFile := filepath.Join(certPath, "localhost+2-key.pem")

	log.Printf("Starting Signaling server at https://localhost:%s", port)
	log.Printf("Using certificates: %s, %s", certFile, keyFile)

	// Запуск сервера с TLS
	if err := router.RunTLS(":"+port, certFile, keyFile); err != nil {
		log.Fatalf("Failed to start Signaling server: %v", err)
	}
}
