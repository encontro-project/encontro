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

	httphandler "encontro/internal/delivery/http"
	"encontro/internal/delivery/websocket"
	"encontro/internal/domain/service"
	"encontro/internal/infrastructure/database"
	postgresRepo "encontro/internal/infrastructure/repository/postgres"
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
	// Инициализация базы данных
	dbConfig := database.NewConfig()
	dbPool, err := database.NewPool(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer dbPool.Close()

	// Инициализация репозиториев
	roomRepo := postgresRepo.NewPostgresRoomRepository(dbPool)
	messageRepo := postgresRepo.NewPostgresMessageRepository(dbPool)
	userSummaryRepo := postgresRepo.NewPostgresUserSummaryRepository(dbPool.GetPool())

	// Инициализация UUID генератора
	uuidGen := service.NewGoogleUUIDGenerator()

	// Инициализация usecase
	messageUseCase := usecase.NewMessageUseCase(messageRepo, uuidGen)
	roomUseCase := usecase.NewRoomUseCase(roomRepo, uuidGen)
	userUseCase := usecase.NewUserSummaryUseCase(userSummaryRepo)

	// Инициализация обработчиков
	wsHandler := websocket.NewHandler(roomUseCase)

	// Настройка маршрутизатора
	router := httphandler.SetupRouter(roomUseCase, messageUseCase, userUseCase, wsHandler)

	// Статические файлы
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	// Загрузка TLS конфигурации
	tlsConfig, err := loadTLSConfig()
	if err != nil {
		log.Fatalf("Failed to load TLS config: %v", err)
	}

	// Создание HTTP сервера
	server := &http.Server{
		Addr:      ":8080",
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	// Канал для получения сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Printf("Starting server at https://localhost:8080")
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
