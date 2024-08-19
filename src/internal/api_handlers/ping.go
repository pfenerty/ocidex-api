package api_handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Ping godoc
// @Summary Ping the server
// @Description Ping the server to check if its running
// @Tags ping
// @Produce  json
// @Success 200 {object} gin.H
// @Router /ping [get]
func ping(ctx *gin.Context, logger *logrus.Logger) {
	logger.Debug("ping request received")
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}
