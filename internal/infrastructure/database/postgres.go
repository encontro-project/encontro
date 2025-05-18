package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config содержит конфигурацию подключения к PostgreSQL
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConfig создает новую конфигурацию, используя переменные окружения или завершает работу с ошибкой
func NewConfig() *Config {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "77.221.159.137"
	}
	portStr := os.Getenv("DB_PORT")
	port := 5432
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		panic("DB_USER не задан в переменных окружения!")
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		panic("DB_PASSWORD не задан в переменных окружения!")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "encontro"
	}
	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "require"
	}
	return &Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}
}

// DSN возвращает строку подключения к PostgreSQL
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// Pool представляет пул соединений с PostgreSQL
type Pool struct {
	pool *pgxpool.Pool
}

// NewPool создает новый пул соединений
func NewPool(ctx context.Context, cfg *Config) (*Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	// Настройка пула соединений
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пула: %w", err)
	}

	// Проверка подключения
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	return &Pool{pool: pool}, nil
}

// Close закрывает пул соединений
func (p *Pool) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

// WithTransaction выполняет функцию в транзакции
func (p *Pool) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := fn(ctx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return nil
}

// GetPool возвращает пул соединений
func (p *Pool) GetPool() *pgxpool.Pool {
	return p.pool
}
