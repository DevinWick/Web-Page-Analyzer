package main

import (
	"github.com/devinwick/web-page-analyzer/handlers"
	LOGGER "github.com/devinwick/web-page-analyzer/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

func init() {
	logger = LOGGER.Log.WithField("pkg", "main")
}

func main() {

	r := gin.Default()
	port := ":8080"

	r.LoadHTMLFiles("pages/index.html", "pages/results.html")
	r.Static("/static", "./static")

	r.GET("/", handlers.IndexPathHandler)
	r.POST("/analyze", handlers.AnalyzeHandler)

	logger.Info("Server starting on port ", port)

	err := r.Run(port)
	if err != nil {
		logger.Fatal("Server start Failed", err)
	}
}
