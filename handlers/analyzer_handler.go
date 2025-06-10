package handlers

import (
	"github.com/devinwick/web-page-analyzer/model"
	"github.com/devinwick/web-page-analyzer/service"
	"github.com/gin-gonic/gin"
	"net/http"
	URL "net/url"
	"strings"
)

func IndexPathHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func AnalyzeHandler(c *gin.Context) {
	rawUrl := c.PostForm("url")

	url := strings.TrimSpace(rawUrl)

	//validate url
	if url == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": "URL is required",
		})
		return
	}

	_, err := URL.ParseRequestURI(url)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": "Invalid URL .",
		})
		return
	}

	//analyze web page
	result, err := service.AnalyzeWebPage(url)
	if err != nil {
		errorResult := model.AnalysisResult{
			URL:        url,
			Error:      err.Error(),
			StatusCode: result.StatusCode,
		}
		c.HTML(http.StatusOK, "results.html", gin.H{
			"result": errorResult,
		})
		return
	}

	c.HTML(http.StatusOK, "results.html", gin.H{
		"result": result,
	})
}
