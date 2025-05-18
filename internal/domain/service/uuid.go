package service

import (
	"github.com/google/uuid"
)

// UUIDGenerator определяет интерфейс для генерации UUID
type UUIDGenerator interface {
	Generate() string
}

// GoogleUUIDGenerator реализует UUIDGenerator используя google/uuid
type GoogleUUIDGenerator struct{}

// NewGoogleUUIDGenerator создает новый экземпляр GoogleUUIDGenerator
func NewGoogleUUIDGenerator() *GoogleUUIDGenerator {
	return &GoogleUUIDGenerator{}
}

// Generate генерирует новый UUID v4
func (g *GoogleUUIDGenerator) Generate() string {
	return uuid.New().String()
}
