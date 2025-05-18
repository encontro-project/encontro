package entity

// PaginationParams содержит параметры пагинации
type PaginationParams struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// PaginatedResponse представляет пагинированный ответ
type PaginatedResponse[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// NewPaginatedResponse создает новый пагинированный ответ
func NewPaginatedResponse[T any](items []T, total int64, params PaginationParams) PaginatedResponse[T] {
	totalPages := (int(total) + params.PageSize - 1) / params.PageSize
	return PaginatedResponse[T]{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}
}
