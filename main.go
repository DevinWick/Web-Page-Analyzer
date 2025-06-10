package main

import (
	"github.com/devinwick/web-page-analyzer/handlers"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	//"web-page-analyzer/handlers"
)

const PORT string = ":8080"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	router := gin.Default()

	//setup templates
	router.LoadHTMLFiles("pages/index.html", "pages/results.html")

	//serve static files
	router.Static("/static", "./static")

	//Configure Routes
	router.GET("/", handlers.IndexPathHandler)

	router.POST("/analyze", handlers.AnalyzeHandler)

	//start server
	err := router.Run(PORT)
	if err != nil {
		logger.Error("Server start Failed", err)
	}
	logger.Info("Server started on port %s", PORT)
}
