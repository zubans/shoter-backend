package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"shuter-go/internal/dto"
	"shuter-go/internal/services"
	"shuter-go/pkg/logger"
)

type UserHandler struct {
	PlayerService *services.PlayerService
}

func New(userService *services.PlayerService) *UserHandler {
	return &UserHandler{PlayerService: userService}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CredentialsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("invalid request", zap.Error(err), zap.String("PlayerId", req.PlayerID), zap.Any("BODY", c.Request.Body))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.PlayerService.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}

	c.Status(http.StatusOK)
}
