package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/pfenerty/ocidex-api/cmd/docs"
	"github.com/pfenerty/ocidex-api/internal/api-handlers"
	"github.com/pfenerty/ocidex-api/middleware"
	"github.com/sirupsen/logrus"
	"os"
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

	r := gin.Default()
	r.Use(middleware.LogrusLogger(logger))
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		logger.Fatal(err)
		return
	}
	api_handlers.RegisterRoutes(r, logger)

	err = r.Run()
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}
}
