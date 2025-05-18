package errors

import "fmt"

// Общие ошибки приложения
var (
	// ErrNotFound возвращается, когда сущность не найдена
	ErrNotFound = func(entity string, id interface{}) error {
		return fmt.Errorf("%s not found: %v", entity, id)
	}

	// ErrInvalidInput возвращается при неверных входных данных
	ErrInvalidInput = func(entity string, reason string) error {
		return fmt.Errorf("invalid %s: %s", entity, reason)
	}

	// ErrInternal возвращается при внутренних ошибках сервера
	ErrInternal = func(operation string, err error) error {
		return fmt.Errorf("internal error during %s: %w", operation, err)
	}
)
