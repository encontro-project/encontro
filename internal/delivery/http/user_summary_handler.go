package http

import (
	"net/http"

	"encontro/internal/delivery/http/middleware"
	"encontro/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserSummaryHandler struct {
	userUseCase *usecase.UserSummaryUseCase
}

func NewUserSummaryHandler(userUC *usecase.UserSummaryUseCase) *UserSummaryHandler {
	return &UserSummaryHandler{userUseCase: userUC}
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
