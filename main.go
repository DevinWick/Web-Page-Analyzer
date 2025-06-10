package main

import (
	"github.com/devinwick/web-page-analyzer/handlers"
	LOGGER "github.com/devinwick/web-page-analyzer/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const PORT string = ":8080"

var logger *logrus.Entry

func init() {
	logger = LOGGER.Log.WithField("pkg", "main")
}

func main() {
	router := gin.Default()

	//setup templates
	router.LoadHTMLFiles("pages/index.html", "pages/results.html")

	//serve static files
	router.Static("/static", "./static")

	//Configure Routes
	router.GET("/", handlers.IndexPathHandler)

	router.POST("/analyze", handlers.AnalyzeHandler)

	//start server
	logger.Info("Server starting on port ", PORT)
	err := router.Run(PORT)
	if err != nil {
		logger.Error("Server start Failed", err)
	}
}
