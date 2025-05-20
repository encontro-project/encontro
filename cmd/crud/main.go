package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	httphandler "encontro/internal/delivery/http"
	middleware "encontro/internal/delivery/http/middleware"
	repoiface "encontro/internal/domain/repository"
	"encontro/internal/domain/service"
	"encontro/internal/infrastructure/database"
	"encontro/internal/infrastructure/repository"
	"encontro/internal/infrastructure/server"
	"encontro/internal/usecase"
)

func main() {
	// Получаем порт из переменных окружения
	port := os.Getenv("CRUD_PORT")
	if port == "" {
		port = "8080" // порт по умолчанию для CRUD
	}

	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Определяем режим работы
	useInMemory := os.Getenv("USE_INMEMORY") == "1"
	initTestData := os.Getenv("INIT_TEST_DATA") == "1"

	var roomRepo repoiface.RoomRepository
	var messageRepo repoiface.MessageRepository
	var userSummaryRepo repoiface.UserSummaryRepository

	if useInMemory {
		// Режим разработки: in-memory хранилище без БД
		log.Println("[INFO] Запуск в режиме разработки: in-memory хранилище")
		if initTestData {
			log.Println("[INFO] Инициализация тестовыми данными")
		}

		// Создаем конфигурацию для in-memory режима
		inMemoryCfg := repository.NewInMemoryConfig(initTestData)
		roomRepo = repository.NewInMemoryRoomRepository(inMemoryCfg)
		messageRepo = repository.NewInMemoryMessageRepository(inMemoryCfg)
		userSummaryRepo = repository.NewInMemoryUserSummaryRepository(inMemoryCfg)
	} else {
		// Режим production: PostgreSQL
		log.Println("[INFO] Запуск в режиме production: PostgreSQL")

		// Проверяем наличие необходимых переменных окружения для БД
		requiredEnvVars := []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_NAME"}
		for _, envVar := range requiredEnvVars {
			if os.Getenv(envVar) == "" {
				log.Fatalf("[ERROR] Необходимо указать переменную окружения %s для работы с PostgreSQL", envVar)
			}
		}

		// Инициализация подключения к БД
		dbConfig := database.NewConfig()
		dbPool, err := database.NewPool(ctx, dbConfig)
		if err != nil {
			log.Fatalf("[ERROR] Ошибка подключения к БД: %v", err)
		}
		defer dbPool.Close()

		roomRepo = repository.NewPostgresRoomRepository(dbPool)
		messageRepo = repository.NewPostgresMessageRepository(dbPool)
		userSummaryRepo = repository.NewPostgresUserSummaryRepository(dbPool.GetPool())
	}

	// Инициализация UUID генератора
	uuidGen := service.NewGoogleUUIDGenerator()

	// Инициализация use cases
	roomUseCase := usecase.NewRoomUseCase(roomRepo, uuidGen)
	messageUseCase := usecase.NewMessageUseCase(messageRepo, uuidGen)
	userSummaryUseCase := usecase.NewUserSummaryUseCase(userSummaryRepo)

	// Инициализация обработчиков
	roomHandler := httphandler.NewRoomHandler(roomUseCase)
	messageHandler := httphandler.NewMessageHandler(messageUseCase)
	userSummaryHandler := httphandler.NewUserSummaryHandler(userSummaryUseCase)

	// Настройка маршрутизатора с конфигурацией по умолчанию
	router := server.NewRouter(nil)

	// API для комнат
	router.POST("/rooms", roomHandler.CreateRoom)
	router.GET("/rooms", roomHandler.ListRooms)
	router.GET("/rooms/:id", roomHandler.GetRoom)
	router.DELETE("/rooms/:id", roomHandler.DeleteRoom)

	// API для сообщений
	router.POST("/rooms/:id/messages", messageHandler.CreateMessage)
	router.GET("/rooms/:id/messages", messageHandler.GetMessages)
	router.GET("/rooms/:id/messages/:messageId", messageHandler.GetMessage)
	router.DELETE("/rooms/:id/messages/:messageId", messageHandler.DeleteMessage)

	router.GET("/api/user/:id", middleware.UserIDMiddleware(), userSummaryHandler.GetUserSummary)

	// Канал для получения сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Printf("Starting CRUD server at http://localhost:%s", port)
		if err := router.Run(":" + port); err != nil {
			log.Printf("Ошибка запуска сервера: %v", err)
			cancel()
		}
	}()

	// Ожидание сигнала завершения
	<-sigChan
	log.Println("Получен сигнал завершения, останавливаем сервер...")
	cancel()
}
