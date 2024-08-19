package api_handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, logger *logrus.Logger) {
	r.GET("/ping", func(ctx *gin.Context) { ping(ctx, logger) })
	r.POST("/image", func(ctx *gin.Context) { getImage(ctx, logger) })
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
