package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/pfenerty/ocidex-api/cmd/docs"
	"github.com/pfenerty/ocidex-api/internal/api-handlers"
	"github.com/pfenerty/ocidex-api/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

// @title OCIDex API
// @version 0.1
// @description OCI metadata API

// @contact.name Patrick Fenerty
// @contact.email patrick@fenerty.me

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	if len(allowedOrigins) == 0 {
		allowedOrigins = append(allowedOrigins, "*")
		logger.Warn("No allowed origin set. Using default.")
	}

	ginMode := os.Getenv("GIN_MODE")
	gin.SetMode(ginMode)

	r := gin.Default()
	r.Use(middleware.LogrusLogger(logger))
	r.Use(cors.New(cors.Config{
		AllowMethods: []string{http.MethodGet, http.MethodPost},
		AllowOrigins: allowedOrigins,
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))
	api_handlers.RegisterRoutes(r, logger)

	err := r.Run()
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}
}
