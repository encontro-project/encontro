package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// ServerConfig представляет конфигурацию HTTP-сервера
type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// NewDefaultServerConfig возвращает конфигурацию сервера по умолчанию
func NewDefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:            "8080",
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 5 * time.Second,
	}
}

// StartServer запускает HTTP-сервер с указанной конфигурацией
func StartServer(cfg *ServerConfig, router chi.Router) error {
	if cfg == nil {
		cfg = NewDefaultServerConfig()
	}

	// Создаем HTTP-сервер
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Запуск сервера на http://localhost:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	// Ожидаем сигнал завершения
	<-ctx.Done()

	// Graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("ошибка при остановке сервера: %v", err)
	}

	return nil
}

// NewChiRouter создает новый роутер с базовыми middleware
func NewChiRouter() chi.Router {
	r := chi.NewRouter()

	// Базовые middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Heartbeat("/ping"))

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	return r
}
