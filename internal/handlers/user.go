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
	userService *services.UserService
}

func New(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CredentialsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("invalid request", zap.Error(err), zap.String("PlayerId", req.PlayerID), zap.Any("BODY", c.Request.Body))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
}
