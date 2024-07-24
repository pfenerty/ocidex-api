package api_handlers

import (
	"github.com/gin-gonic/gin"
	registryMetadata "github.com/pfenerty/ocidex-api/internal/registry-metadata"
	"github.com/sirupsen/logrus"
	"net/http"
)

// getImage godoc
// @Summary Get image metadata for an image reference
// @Description Get image manifest and config
// @Tags images
// @Accept  json
// @Produce  json
// @Param   ImageRequestQuery     body    registryMetadata.ImageRequestQuery     true        "Image Request Query"
// @Success 200 {object} registryMetadata.ImageResponseData
// @Failure 400 {object} gin.H
// @Router /image [post]
func getImage(ctx *gin.Context, logger *logrus.Logger) {
	logger.Debug("get image request received")

	var requestParams registryMetadata.ImageRequestQuery
	if err := ctx.ShouldBindJSON(&requestParams); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to bind request body to docker reference")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry := logger.WithFields(logrus.Fields{
		"reference": requestParams.Reference,
	})
	entry.Debug("getting image metadata")

	if requestParams.Architecture == "" {
		requestParams.Architecture = "amd64"
	}

	if requestParams.OS == "" {
		requestParams.OS = "linux"
	}

	metadata, err := registryMetadata.ConstructRegistryMetadata(entry, requestParams)
	if err != nil {
		entry.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to construct metadata")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, metadata)
}
