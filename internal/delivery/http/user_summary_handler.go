package http

import (
	"net/http"

	"encontro/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

// UserSummaryHandler обрабатывает HTTP-запросы для пользователей
type UserSummaryHandler struct {
	userUseCase UserSummaryUseCaseInterface
}

// NewUserSummaryHandler создает новый экземпляр UserSummaryHandler
func NewUserSummaryHandler(userUseCase UserSummaryUseCaseInterface) *UserSummaryHandler {
	return &UserSummaryHandler{userUseCase: userUseCase}
}

func (h *UserSummaryHandler) GetUserSummary(c *gin.Context) {
	userID, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user id missing in context"})
		return
	}

	summary, err := h.userUseCase.GetUserSummary(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
