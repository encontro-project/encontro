package repository

import "sync"

// BaseRepository предоставляет базовую функциональность для всех репозиториев
type BaseRepository struct {
	mu sync.RWMutex
}

// Lock блокирует репозиторий для записи
func (r *BaseRepository) Lock() {
	r.mu.Lock()
}

// Unlock разблокирует репозиторий после записи
func (r *BaseRepository) Unlock() {
	r.mu.Unlock()
}

// RLock блокирует репозиторий для чтения
func (r *BaseRepository) RLock() {
	r.mu.RLock()
}

// RUnlock разблокирует репозиторий после чтения
func (r *BaseRepository) RUnlock() {
	r.mu.RUnlock()
}

// WithWriteLock выполняет операцию с блокировкой на запись
func (r *BaseRepository) WithWriteLock(fn func()) {
	r.Lock()
	defer r.Unlock()
	fn()
}

// WithReadLock выполняет операцию с блокировкой на чтение
func (r *BaseRepository) WithReadLock(fn func()) {
	r.RLock()
	defer r.RUnlock()
	fn()
}
